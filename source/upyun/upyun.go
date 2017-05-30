// +-------------------------------------------------------------------------
// | Copyright (C) 2017 Yunify, Inc.
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

package upyun

import (
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/frostyplanet/logrus"
	"github.com/upyun/go-sdk/upyun"
	"github.com/yunify/qscamel/record"
)

const upyunEndpoint = "upaiyun.com"

// SourceUpyun is upyun bucket type source
type SourceUpyun struct {
	BucketURL  string
	bucketName string
	up         *upyun.UpYun

	lastModifiedTimes map[string]time.Time
	objsChan          chan *upyun.FileInfo
	marker            string
}

// NewSourceUpyun creates an instance of SourceUpyun
func NewSourceUpyun(bucketName, zone, accessKeyID, secretAccessKey string) (*SourceUpyun, error) {
	if zone == "" || accessKeyID == "" || secretAccessKey == "" {
		return &SourceUpyun{}, errors.New("Miss zone, accessKeyID or secretAccessKey for upyun")
	}

	up := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   bucketName,
		Operator: accessKeyID,
		Password: secretAccessKey,
	})

	return &SourceUpyun{
		BucketURL:         "http://" + bucketName + "." + zone + "." + upyunEndpoint,
		bucketName:        bucketName,
		up:                up,
		lastModifiedTimes: make(map[string]time.Time),
		marker:            "",
	}, nil
}

// GetSourceSites implements MigrateSource.GetSourceSites
func (source *SourceUpyun) GetSourceSites(
	threadNum int, logger *log.Logger, recorder *record.Recorder,
) (sourceSites []string, objectNames []string, skipped []string, done bool, err error) {
	sourceSites, objectNames, skipped = []string{}, []string{}, []string{}
	err = nil

	if source.marker == "" {
		source.objsChan = make(chan *upyun.FileInfo, 10)
		// upyun golang sdk does not suport to get part of all the objects to do advanced treatment.
		// Inside the "List" interface, it Summarize all objects got from diffent http resp
		// So we store all the objects only at first time.
		// Then we use the stored objects to do partial traversal regarding the number of concurrent tasks
		getObjCfg := upyun.GetObjectsConfig{
			Path:         "/",
			ObjectsChan:  source.objsChan,
			MaxListLevel: -1,
		}
		err = source.up.List(&getObjCfg)

		if err != nil {
			return sourceSites, objectNames, skipped, done, err
		}

		// If the implementation of the “List” interface is changed in the later version ，the marker will be a cursor.
		// But at current version, it is just a flag which indicate that the information of objects has been acquired.
		source.marker = "The information of objects has been acquired."
	}

	// If we got enough sourceSites as mouch as threadNum, return and begin to fetch these objects concurrently
	for i := 0; i < threadNum; {

		if len(source.objsChan) == 0 {
			// If there is no more objects to fetch, the entire task is fished.
			done = true
			return sourceSites, objectNames, skipped, done, err
		}

		for object := range source.objsChan {

			if object.IsDir {
				logger.Infof("Skip dir object: %s", object.Name)
				continue
			}

			sourceSite := strings.Join([]string{source.BucketURL, object.Name}, "/")
			if recorder.IsExist(sourceSite) {
				logger.Infof("Skip completed source site: %s", sourceSite)
				skipped = append(skipped, sourceSite)
				continue
			}

			sourceSites = append(sourceSites, sourceSite)
			objectNames = append(objectNames, object.Name)
			source.lastModifiedTimes[sourceSite] = object.Time

			i++
		}
	}

	return sourceSites, objectNames, skipped, done, err
}

// GetSourceSiteInfo implements MigrateSource.GetSourceSites
func (source *SourceUpyun) GetSourceSiteInfo(sourceSite string) (lastModified time.Time, err error) {
	if lastModified, ok := source.lastModifiedTimes[sourceSite]; ok {
		return lastModified, nil
	}
	return time.Time{}, fmt.Errorf(
		"SourceUpyun hasn't gotten source site %s from bucket %s", sourceSite, source.bucketName,
	)
}
