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

package qiniu

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	log "github.com/frostyplanet/logrus"
	"github.com/qiniu/api.v6/auth/digest"
	"github.com/qiniu/api.v6/conf"
	"github.com/qiniu/api.v6/rs"
	"github.com/qiniu/api.v6/rsf"

	"github.com/yunify/qscamel/record"
)

// SourceQiniu is qiniu bucket type source
type SourceQiniu struct {
	BucketURL  string
	bucketName string
	// Key:sourceSite Value:LastModifiedTime
	lastModifiedTimes map[string]time.Time
}

var errMissCredential = errors.New("Miss accessKeyID or secretAccessKey for qiniu")

// NewSourceQiniu creates an instance of SourceQiniu
func NewSourceQiniu(bucketName, accessKeyID, secretAccessKey string) (SourceQiniu, error) {
	if accessKeyID == "" || secretAccessKey == "" {
		return SourceQiniu{}, errMissCredential
	}
	conf.ACCESS_KEY = accessKeyID
	conf.SECRET_KEY = secretAccessKey
	bucketURL, err := getBucketURL(bucketName, accessKeyID, secretAccessKey)
	if err != nil {
		return SourceQiniu{}, err
	}
	return SourceQiniu{
		BucketURL:         "http://" + bucketURL,
		bucketName:        bucketName,
		lastModifiedTimes: make(map[string]time.Time),
	}, nil
}

func getBucketURL(bucket, accessKeyID, secretAccessKey string) (url string, err error) {
	domains := make([]string, 0)
	mac := &digest.Mac{AccessKey: accessKeyID, SecretKey: []byte(secretAccessKey)}
	client := rs.New(mac)
	getDomainsURL := "http://api.qiniu.com/v6/domain/list"
	postData := map[string][]string{
		"tbl": []string{bucket},
	}
	err = client.Conn.CallWithForm(nil, &domains, getDomainsURL, postData)
	if err == nil {
		url = domains[0]
	}
	return
}

// GetSourceSites implements MigrateSource.GetSourceSites
func (source SourceQiniu) GetSourceSites(
	threadNum int, logger *log.Logger, recorder *record.Recorder,
) (sourceSites []string, objectNames []string, skipped []string, done bool, err error) {
	sourceSites, objectNames, skipped = []string{}, []string{}, []string{}
	marker := ""
	for i := 0; i < threadNum; {
		if done {
			return
		}
		limit := threadNum - i
		client := rsf.New(nil)
		objects, nextMarker, err := client.ListPrefix(nil, source.bucketName, "", marker, limit)
		if err != nil {
			if err == io.EOF {
				done = true
			} else {
				return sourceSites, objectNames, skipped, done, err
			}
		}
		marker = nextMarker

		for _, object := range objects {
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
			source.lastModifiedTimes[sourceSite] = time.Unix(0, object.PutTime*100)
		}
	}
	return
}

// GetSourceSiteInfo implements MigrateSource.GetSourceSites
func (source SourceQiniu) GetSourceSiteInfo(sourceSite string) (lastModified time.Time, err error) {
	if lastModified, ok := source.lastModifiedTimes[sourceSite]; ok {
		return lastModified, nil
	}
	return time.Time{}, fmt.Errorf(
		"SourceQiniu hasn't gotten source site %s from bucket %s", sourceSite, source.bucketName,
	)
}
