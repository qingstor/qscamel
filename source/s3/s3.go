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

package s3

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/frostyplanet/logrus"

	"github.com/yunify/qscamel/record"
)

const s3Endpoint = ".s3.amazonaws.com"

// SourceS3 is s3 bucket type source
type SourceS3 struct {
	BucketURL  string
	bucketName string
	service    *s3.S3
	// Key:sourceSite Value:LastModifiedTime
	lastModifiedTimes map[string]time.Time
}

var errMissCredential = errors.New("Miss zone, accessKeyID or secretAccessKey for s3")

// NewSourceS3 creates an instance of SourceS3
func NewSourceS3(bucketName, zone, accessKeyID, secretAccessKey string) (SourceS3, error) {
	if zone == "" || accessKeyID == "" || secretAccessKey == "" {
		return SourceS3{}, errMissCredential
	}
	cfg := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			accessKeyID, secretAccessKey, "",
		),
		Region: &zone,
	}
	sess, err := session.NewSession(cfg)
	if err != nil {
		return SourceS3{}, fmt.Errorf("SourceS3 failed to create session, %v", err)
	}
	service := s3.New(sess)
	return SourceS3{
		BucketURL:         "http://" + bucketName + s3Endpoint,
		bucketName:        bucketName,
		service:           service,
		lastModifiedTimes: make(map[string]time.Time),
	}, nil
}

// GetSourceSites implements MigrateSource.GetSourceSites
func (source SourceS3) GetSourceSites(
	threadNum int, logger *log.Logger, recorder *record.Recorder,
) (sourceSites []string, objectNames []string, skipped []string, done bool, err error) {
	sourceSites, objectNames, skipped = []string{}, []string{}, []string{}
	marker := ""
	for i := 0; i < threadNum; {
		if done {
			return
		}
		maxKeys := int64(threadNum - i)
		input := &s3.ListObjectsV2Input{
			Bucket:     &source.bucketName,
			MaxKeys:    &maxKeys,
			StartAfter: &marker,
		}
		result, err := source.service.ListObjectsV2(input)
		if err != nil {
			return sourceSites, objectNames, skipped, done, err
		}
		// IsTruncated is true: there are more keys in the bucket
		if !*result.IsTruncated {
			done = true
		}
		if result.ContinuationToken != nil {
			marker = *result.ContinuationToken
		}

		for _, object := range result.Contents {
			sourceSite := strings.Join([]string{source.BucketURL, *object.Key}, "/")
			if recorder.IsExist(sourceSite) {
				logger.Infof(
					"Skip completed source site: %s", sourceSite,
				)
				skipped = append(skipped, sourceSite)
				continue
			}
			i++
			sourceSites = append(sourceSites, sourceSite)
			objectNames = append(objectNames, *object.Key)
			source.lastModifiedTimes[sourceSite] = *object.LastModified
		}
	}
	return
}

// GetSourceSiteInfo implements MigrateSource.GetSourceSites
func (source SourceS3) GetSourceSiteInfo(sourceSite string) (lastModified time.Time, err error) {
	if lastModified, ok := source.lastModifiedTimes[sourceSite]; ok {
		return lastModified, nil
	}
	return time.Time{}, fmt.Errorf(
		"SourceS3 hasn't gotten source site %s from bucket %s", sourceSite, source.bucketName,
	)
}
