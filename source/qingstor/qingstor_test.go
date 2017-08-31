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
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	qs "github.com/yunify/qingstor-sdk-go/service"
	"github.com/yunify/qscamel/record"
	"github.com/yunify/qscamel/utils"
)

const (
	envBucketName = "QSCAMEL_TEST_SOURCEQINGSTOR_BUCKETNAME"
	envZone       = "QSCAMEL_TEST_SOURCEQINGSTOR_ZONE"
	envAccessKey  = "QSCAMEL_TEST_SOURCEQINGSTOR_ACCESS_KEY_ID"
	envSecretKey  = "QSCAMEL_TEST_SOURCEQINGSTOR_SECRET_ACCESS_KEY"
)

type testCase struct {
	SourceObjects       []string
	ThreadNum           int
	RecordObjects       []string
	ExpectedObjectNames []string
	ExpectedSkipped     []string
	ExpectedDone        bool
}

var sourceQingstorTests = []testCase{
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

func TestSourceQingstor_GetSourceSites(t *testing.T) {
	bucketName := os.Getenv(envBucketName)
	zone := os.Getenv(envZone)
	accessKeyID := os.Getenv(envAccessKey)
	secretAccessKey := os.Getenv(envSecretKey)
	if bucketName == "" || zone == "" || accessKeyID == "" || secretAccessKey == "" {
		fmt.Println("Miss environment variables for SourceQingstor test")
		t.SkipNow()
	}

	source, err := NewSourceQingstor(bucketName, zone, accessKeyID, secretAccessKey)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	logger := utils.GetLogger()
	// Clear test bucket
	clearBucket(source)
	bucket, _ := source.service.Bucket(source.bucketName, source.zone)
	for _, test := range sourceQingstorTests {
		for _, object := range test.SourceObjects {
			file, _ := os.Create(object)
			_, err := bucket.PutObject(object, &qs.PutObjectInput{Body: file})
			assert.Nil(t, err)
		}
		recordFile, err := record.GetRecordFile("test_get_source_sites")
		assert.Nil(t, err)
		if test.RecordObjects != nil {
			recordSourceSites := getSourceSites(source, test.RecordObjects)
			for _, record := range recordSourceSites {
				recordFile.WriteString(fmt.Sprintf("%s\n", record))
			}
			recordFile.Seek(0, 0)
		}
		recorder := record.NewRecorder(recordFile)

		sourceSites, names, skipped, done, err := source.GetSourceSites(test.ThreadNum, logger, recorder)

		assert.Nil(t, err)

		expectedSourceSites := getSourceSites(source, test.ExpectedObjectNames)
		assert.Equal(t, expectedSourceSites, sourceSites)
		assert.Equal(t, test.ExpectedObjectNames, names)
		expectedSkipped := getSourceSites(source, test.ExpectedSkipped)
		assert.Equal(t, expectedSkipped, skipped)
		assert.Equal(t, test.ExpectedDone, done)

		clearBucket(source)
		clearSourceQingstor(source)
		recorder.Clear()
	}
}

func getSourceSites(source *SourceQingstor, objectNames []string) []string {
	sourceSites := []string{}
	for _, object := range objectNames {
		sourceSite := fmt.Sprintf("%s/%s", source.BucketURL, object)
		httpRequest, _ := http.NewRequest("GET", sourceSite, nil)
		expires := time.Now().Unix() + 60*60*24
		signature, _ := QSSigner.BuildQuerySignature(httpRequest, (int)(expires))
		sourceSite = fmt.Sprintf("%s?%s", sourceSite, signature)
		sourceSites = append(
			sourceSites, sourceSite,
		)
	}
	return sourceSites
}

func clearBucket(source *SourceQingstor) {
	bucket, _ := source.service.Bucket(source.bucketName, source.zone)
	listObjects, _ := bucket.ListObjects(nil)
	for _, object := range listObjects.Keys {
		bucket.DeleteObject(*object.Key)
	}
}

func clearSourceQingstor(source *SourceQingstor) {
	source.lastModifiedTimes = make(map[string]time.Time)
	source.marker = nil
}
