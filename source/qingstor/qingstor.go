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

package qingstor

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/frostyplanet/logrus"

	"github.com/yunify/qingstor-sdk-go/config"
	"github.com/yunify/qingstor-sdk-go/request/signer"
	qs "github.com/yunify/qingstor-sdk-go/service"
	"github.com/yunify/qscamel/record"
)

const (
	qsEndPoint = ".qingstor.com"
	protHeader = "http://"
	slash      = "/"
	reqMethod  = "GET"
)

// QSSigner is the http request signer for QingStor service.
var QSSigner signer.QingStorSigner

// SourceQingstor is qingstor bucket type source.
type SourceQingstor struct {
	BucketURL  string
	bucketName string
	service    *qs.Service
	zone       string
	// Key:sourceSite Value:LastModifiedTime
	lastModifiedTimes map[string]time.Time
	marker            *string
}

// NewSourceQingstor creates an instance of SourceQingstor.
func NewSourceQingstor(bucketName, zone, accessKeyID, secretAccessKey string) (*SourceQingstor, error) {
	if zone == "" || accessKeyID == "" || secretAccessKey == "" {
		return &SourceQingstor{}, errors.New("Miss zone, accessKeyID or secretAccessKey for QingStor")
	}

	conf, err := config.New(accessKeyID, secretAccessKey)
	if err != nil {
		return &SourceQingstor{}, err
	}

	service, err := qs.Init(conf)
	if err != nil {
		return &SourceQingstor{}, err
	}
	QSSigner.AccessKeyID = accessKeyID
	QSSigner.SecretAccessKey = secretAccessKey

	return &SourceQingstor{
		BucketURL:         fmt.Sprintf("%s%s%s%s%s", protHeader, zone, qsEndPoint, slash, bucketName),
		bucketName:        bucketName,
		service:           service,
		zone:              zone,
		lastModifiedTimes: make(map[string]time.Time),
		marker:            nil,
	}, nil
}

// GetSourceSites implements MigrateSource.GetSourceSites
func (source *SourceQingstor) GetSourceSites(threadNum int, logger *log.Logger, recorder *record.Recorder) (sourceSites []string, objectNames []string, skipped []string, done bool, err error) {
	sourceSites, objectNames, skipped = []string{}, []string{}, []string{}

	bucket, err := source.service.Bucket(source.bucketName, source.zone)
	if err != nil {
		return sourceSites, objectNames, skipped, done, err
	}
	for i := 0; i < threadNum; {
		if done {
			return
		}
		limit := threadNum - i
		listObjects, err := bucket.ListObjects(&qs.
			ListObjectsInput{Marker: source.marker, Limit: &limit})
		if err != nil {
			return sourceSites, objectNames, skipped, done, err
		}

		source.marker = listObjects.NextMarker

		for _, object := range listObjects.Keys {
			sourceSite := strings.Join([]string{source.BucketURL, *object.Key}, slash)
			httpRequest, err := http.NewRequest("GET", sourceSite, nil)
			if err != nil {
				return sourceSites, objectNames, skipped, done, err
			}

			expires := time.Now().Unix() + 60*60*24
			signature, err := QSSigner.BuildQuerySignature(httpRequest, (int)(expires))
			if err != nil {
				return sourceSites, objectNames, skipped, done, err
			}
			sourceSite = fmt.Sprintf("%s?%s", sourceSite, signature)

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
			source.lastModifiedTimes[sourceSite] =
				time.Unix(object.Created.Unix(), (int64)(*object.Modified)*100)
		}

		if *source.marker == "" {
			done = true
		}
	}
	return
}

// GetSourceSiteInfo implements MigrateSource.GetSourceSites
func (source *SourceQingstor) GetSourceSiteInfo(sourceSite string) (lastModified time.Time, err error) {
	if lastModified, ok := source.lastModifiedTimes[sourceSite]; ok {
		return lastModified, nil
	}
	return time.Time{}, fmt.Errorf(
		"SourceQingstor hasn't gotten source site %s from bucket %s", sourceSite, source.bucketName,
	)
}
