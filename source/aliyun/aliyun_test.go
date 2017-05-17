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
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qscamel/record"
	"github.com/yunify/qscamel/utils"
)

const (
	envBucketName = "QSCAMEL_TEST_SOURCEALIYUN_BUCKETNAME"
	envZone       = "QSCAMEL_TEST_SOURCEALIYUN_ZONE"
	envAccessKey  = "QSCAMEL_TEST_SOURCEALIYUN_ACCESS_KEY_ID"
	envSecretKey  = "QSCAMEL_TEST_SOURCEALIYUN_SECRET_ACCESS_KEY"
)

type testCase struct {
	SourceObjects       []string
	ThreadNum           int
	RecordObjects       []string
	ExpectedObjectNames []string
	ExpectedSkipped     []string
	ExpectedDone        bool
}

var sourceAliyunTests = []testCase{
	// Normal
	{
		SourceObjects:       []string{"object1", "object2", "object3"},
		ThreadNum:           3,
		ExpectedObjectNames: []string{"object1", "object2", "object3"},
		// IsTruncated is false when limit(maxKeys) is exactly equal to object num.
		ExpectedDone: false,
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
	// TreadNum < total number of objects and TreadNum > number of expected objects
	{
		SourceObjects:       []string{"object1", "object2", "object3"},
		ThreadNum:           2,
		RecordObjects:       []string{"object1", "object2"},
		ExpectedObjectNames: []string{"object3"},
		ExpectedSkipped:     []string{"object1", "object2"},
		ExpectedDone:        true,
	},
}

func TestSourceAliyun_GetSourceSites(t *testing.T) {
	bucketName := os.Getenv(envBucketName)
	zone := os.Getenv(envZone)
	accessKey := os.Getenv(envAccessKey)
	secretKey := os.Getenv(envSecretKey)
	if bucketName == "" || zone == "" || accessKey == "" || secretKey == "" {
		fmt.Println("Miss environment variables for SourceAliyun test")
		t.SkipNow()
	}
	source, err := NewSourceAliyun(bucketName, zone, accessKey, secretKey)
	assert.Nil(t, err)
	logger := utils.GetLogger()

	for _, test := range sourceAliyunTests {
		for _, object := range test.SourceObjects {
			err = source.bucket.PutObject(object, bytes.NewBuffer([]byte("hello world")))
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

		clearBucket(source.bucket)
		clearSourceAliyun(source)
		recorder.Clear()
	}
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

func clearBucket(bucket *oss.Bucket) {
	result, err := bucket.ListObjects()
	if err == nil {
		for _, object := range result.Objects {
			bucket.DeleteObject(object.Key)
		}
	}
}

func clearSourceAliyun(source *SourceAliyun) {
	source.lastModifiedTimes = make(map[string]time.Time)
	source.firstObject = ""
	source.marker = ""
}
