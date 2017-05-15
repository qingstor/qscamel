// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package migrate

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/frostyplanet/logrus"

	"github.com/yunify/qingstor-sdk-go/request/errs"
	"github.com/yunify/qingstor-sdk-go/service"
	"github.com/yunify/qscamel/utils"
)

type fetchResult struct {
	Completed bool
	Source    string
}

func checkIgnoreCondition(context *Context, bucket *service.Bucket, objectName string,
	sourceSite string, failed, skipped *[]string) (ignore bool) {
	if context.IgnoreExisting || context.IgnoreUnmodified {
		objectInfo, err := bucket.HeadObject(
			objectName, &service.HeadObjectInput{},
		)
		if err != nil {
			// If object doesn't exist, still fetch
			qsErr, ok := err.(*errs.QingStorError)
			if !ok || qsErr.StatusCode != http.StatusNotFound {
				context.Logger.Errorf(
					"Error occurs when heading object %s. %s",
					objectName, err.Error(),
				)
				*failed = append(*failed, sourceSite)
				ignore = true
			}
		} else {
			if context.IgnoreExisting {
				// Skip existing object.
				context.Logger.Infof("Skip existing object: %s", objectName)
				*skipped = append(*skipped, sourceSite)
				ignore = true
			} else {
				// IgnoreUnmodified is true, check whether object is the latest.
				sourceLastModified, err := context.Source.GetSourceSiteInfo(sourceSite)
				if err != nil {
					context.Logger.Warnf(
						"Can't get last modified time of source site %s. %v",
						sourceSite, err,
					)
				}
				if objectInfo.LastModified.Local().After(sourceLastModified.Local()) {
					context.Logger.Infof("Skip the latest object: %s", objectName)
					*skipped = append(*skipped, sourceSite)
					ignore = true
				}
			}
		}
	}
	return
}

// Migrate reads source list, executes migration and waits for all migrations done.
// It returns three string slice for completed, failed and skipped situation.
func Migrate(context *Context) ([]string, []string, []string, error) {
	completed, failed, skipped := []string{}, []string{}, []string{}

	service, _ := service.Init(context.QSConfig)
	zone, err := utils.DetermineBucketZone(service, context.QSBucketName)
	if err != nil {
		return completed, failed, skipped, err
	}
	bucket, _ := service.Bucket(context.QSBucketName, zone)
	for {
		sourceSites, objectNames, skippedSourceFiles, endOfSource, err :=
			context.Source.GetSourceSites(context.ThreadNum, context.Logger, context.Recorder)
		if err != nil {
			return completed, failed, skipped, err
		}
		skipped = append(skipped, skippedSourceFiles...)
		var resultChan chan fetchResult
		if len(sourceSites) != 0 {
			resultChan = make(chan fetchResult, len(sourceSites))
		}

		fetchNum := 0
		for i, sourceSite := range sourceSites {
			objectName := objectNames[i]
			// If not overwrite, check whether skipping existing or unmodified objects.
			if !context.Overwrite && checkIgnoreCondition(
				context, bucket, objectName, sourceSite, &failed, &skipped) {
				continue
			}
			if context.DryRun {
				completed = append(completed, sourceSite)
				continue
			}
			fetchNum++
			go fetchObject(
				objectName, sourceSite, bucket, resultChan,
				context.Logger,
			)
		}

		// Wait for completion of this batch
		progressBarTimer := time.NewTimer(time.Second * 2)
		for i := 0; i < fetchNum; {
			select {
			case <-progressBarTimer.C:
				progressBarTimer.Reset(time.Second * 2)
				fmt.Print(">>")
			case result := <-resultChan:
				if result.Completed {
					completed = append(completed, result.Source)
					context.Recorder.Put(result.Source)

				} else {
					failed = append(failed, result.Source)
				}
				i++
				fmt.Printf("\n[ %d/%d of the download tasks in current batch is finished. ]\n", i, fetchNum)
			}
		}
		if endOfSource {
			break
		}
		fmt.Println("[ New batch of fiels begin to download. ]")
	}
	fmt.Println("[ All download taskes are finished. ]")
	context.Recorder.Clear()
	return completed, failed, skipped, nil
}

func fetchObject(objectName string, sourceSite string, bucket *service.Bucket,
	resultChan chan fetchResult, logger *log.Logger) {
	_, err := bucket.PutObject(
		objectName,
		&service.PutObjectInput{XQSFetchSource: sourceSite},
	)
	if err != nil {
		logger.Warnf(
			"Can't fetch object %s. %s",
			objectName,
			err.Error(),
		)
		resultChan <- fetchResult{false, sourceSite}
	} else {
		resultChan <- fetchResult{true, sourceSite}
	}
}
