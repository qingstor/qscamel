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
	"bytes"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qscamel/record"
	"github.com/yunify/qscamel/utils"
)

const (
	envZone      = "QSCAMEL_TEST_SOURCES3_ZONE"
	envAccessKey = "QSCAMEL_TEST_SOURCES3_ACCESS_KEY_ID"
	envSecretKey = "QSCAMEL_TEST_SOURCES3_SECRET_ACCESS_KEY"
)

type testCase struct {
	SourceObjects       []string
	ThreadNum           int
	RecordObjects       []string
	ExpectedObjectNames []string
	ExpectedSkipped     []string
	ExpectedDone        bool
}

var sourceS3Tests = []testCase{
	// Normal
	{
		SourceObjects:       []string{"object1", "object2", "object3"},
		ThreadNum:           3,
		ExpectedObjectNames: []string{"object1", "object2", "object3"},
		ExpectedDone:        true,
	},
	// TreadNum > total number of objects
	{
		SourceObjects:       []string{"object1", "object2", "object3"},
		ThreadNum:           5,
		ExpectedObjectNames: []string{"object1", "object2", "object3"},
		ExpectedDone:        true,
	},
	// Skip recorded objects
	{
		SourceObjects:       []string{"object1", "object2", "object3"},
		ThreadNum:           3,
		RecordObjects:       []string{"object1"},
		ExpectedObjectNames: []string{"object2", "object3"},
		ExpectedSkipped:     []string{"object1"},
		ExpectedDone:        true,
	},
}

func TestSourceS3_GetSourceSites(t *testing.T) {
	zone := os.Getenv(envZone)
	accessKey := os.Getenv(envAccessKey)
	secretKey := os.Getenv(envSecretKey)
	if zone == "" || accessKey == "" || secretKey == "" {
		fmt.Println("Miss environment variables for SourceS3 test")
		t.SkipNow()
	}
	bucketName := "test-s3-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	source, err := NewSourceS3(bucketName, zone, accessKey, secretKey)
	assert.Nil(t, err)
	_, err = source.service.CreateBucket(&awss3.CreateBucketInput{Bucket: &bucketName})
	assert.Nil(t, err)
	logger := utils.GetLogger()
	for _, test := range sourceS3Tests {
		for _, object := range test.SourceObjects {
			input := &awss3.PutObjectInput{
				Bucket: &bucketName,
				Key:    &object,
				Body:   bytes.NewReader([]byte("hello world")),
			}
			_, err := source.service.PutObject(input)
			assert.Nil(t, err)
		}
		recordFile, err := record.GetRecordFile("test_get_source_sites")
		assert.Nil(t, err)
		if test.RecordObjects != nil {
			recordSourceSites := getSourceSites(source.BucketURL, test.RecordObjects)
			for _, record := range recordSourceSites {
				recordFile.WriteString(fmt.Sprintf("%s\n", record))
			}
			recordFile.Seek(0, 0)
		}
		recorder := record.NewRecorder(recordFile)

		sourceSites, names, skipped, done, err := source.GetSourceSites(test.ThreadNum, logger, recorder)

		assert.Nil(t, err)
		expectedSourceSites := getSourceSites(source.BucketURL, test.ExpectedObjectNames)
		assert.Equal(t, expectedSourceSites, sourceSites)
		assert.Equal(t, test.ExpectedObjectNames, names)
		expectedSkipped := getSourceSites(source.BucketURL, test.ExpectedSkipped)
		assert.Equal(t, expectedSkipped, skipped)
		assert.Equal(t, test.ExpectedDone, done)

		clearBucket(source.service, bucketName)
		recorder.Clear()
	}
	source.service.DeleteBucket(&awss3.DeleteBucketInput{Bucket: &bucketName})
}

func getSourceSites(upstreamURL string, objectNames []string) []string {
	sourceSites := []string{}
	for _, object := range objectNames {
		sourceSites = append(
			sourceSites, fmt.Sprintf("%s/%s", upstreamURL, object),
		)
	}
	return sourceSites
}

func clearBucket(service *awss3.S3, bucketName string) {
	resp, err := service.ListObjectsV2(&awss3.ListObjectsV2Input{Bucket: &bucketName})
	if err == nil {
		for _, object := range resp.Contents {
			input := &awss3.DeleteObjectInput{
				Bucket: &bucketName,
				Key:    object.Key,
			}
			service.DeleteObject(input)
		}
	}
}
