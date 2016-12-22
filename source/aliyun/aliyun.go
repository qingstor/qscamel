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

package aliyun

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	log "github.com/frostyplanet/logrus"

	"github.com/yunify/qscamel/record"
)

const aliyunEndpoint = ".aliyuncs.com"

// SourceAliyun is aliyun bucket type source
type SourceAliyun struct {
	BucketURL string
	bucket    *oss.Bucket
	// Key:sourceSite Value:LastModifiedTime
	lastModifiedTimes map[string]time.Time
	firstObject       string
}

var errMissCredential = errors.New("Miss zone, accessKeyID or secretAccessKey for aliyun")

// NewSourceAliyun creates an instance of SourceAliyun
func NewSourceAliyun(bucketName, zone, accessKeyID, secretAccessKey string) (SourceAliyun, error) {
	if zone == "" || accessKeyID == "" || secretAccessKey == "" {
		return SourceAliyun{}, errMissCredential
	}
	client, err := oss.New("http://"+zone+aliyunEndpoint, accessKeyID, secretAccessKey)
	if err != nil {
		return SourceAliyun{}, err
	}
	bucket, _ := client.Bucket(bucketName)
	return SourceAliyun{
		BucketURL:         "http://" + bucketName + "." + zone + aliyunEndpoint,
		bucket:            bucket,
		lastModifiedTimes: make(map[string]time.Time),
	}, nil
}

// GetSourceSites implements MigrateSource.GetSourceSites
func (source SourceAliyun) GetSourceSites(
	threadNum int, logger *log.Logger, recorder *record.Recorder,
) (sourceSites []string, objectNames []string, skipped []string, done bool, err error) {
	sourceSites, objectNames, skipped = []string{}, []string{}, []string{}
	marker := ""
	for i := 0; i < threadNum; {
		maxKeys := threadNum - i
		result, err := source.bucket.ListObjects(oss.MaxKeys(maxKeys), oss.Marker(marker))
		if err != nil {
			return sourceSites, objectNames, skipped, done, err
		}
		marker = result.NextMarker

		for _, object := range result.Objects {
			if source.firstObject == "" {
				source.firstObject = object.Key
			} else if source.firstObject == object.Key {
				done = true
				return sourceSites, objectNames, skipped, done, err
			}
			sourceSite := strings.Join([]string{source.BucketURL, object.Key}, "/")
			if recorder.IsExist(sourceSite) {
				logger.Infof(
					"Skip completed source site: %s", sourceSite,
				)
				skipped = append(skipped, sourceSite)
				continue
			}
			i++
			sourceSites = append(sourceSites, sourceSite)
			objectNames = append(objectNames, object.Key)
			source.lastModifiedTimes[sourceSite] = object.LastModified
		}
	}
	return
}

// GetSourceSiteInfo implements MigrateSource.GetSourceSites
func (source SourceAliyun) GetSourceSiteInfo(sourceSite string) (lastModified time.Time, err error) {
	if lastModified, ok := source.lastModifiedTimes[sourceSite]; ok {
		return lastModified, nil
	}
	return time.Time{}, fmt.Errorf(
		"SourceAliyun hasn't gotten source site %s from bucket %s", sourceSite, source.bucket.BucketName,
	)
}
