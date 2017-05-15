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
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/qiniu/api.v6/auth/digest"
	qnio "github.com/qiniu/api.v6/io"
	"github.com/qiniu/api.v6/rs"
	"github.com/qiniu/api.v6/rsf"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qscamel/record"
	"github.com/yunify/qscamel/utils"
)

const (
	envBucketName = "QSCAMEL_TEST_SOURCEQINIU_BUCKETNAME"
	envAccessKey  = "QSCAMEL_TEST_SOURCEQINIU_ACCESS_KEY_ID"
	envSecretKey  = "QSCAMEL_TEST_SOURCEQINIU_SECRET_ACCESS_KEY"
)

type testCase struct {
	SourceObjects       []string
	ThreadNum           int
	RecordObjects       []string
	ExpectedObjectNames []string
	ExpectedSkipped     []string
	ExpectedDone        bool
}

var sourceQiniuTests = []testCase{
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

func TestSourceQiniu_GetSourceSites(t *testing.T) {
	bucketName := os.Getenv(envBucketName)
	accessKeyID := os.Getenv(envAccessKey)
	secretAccessKey := os.Getenv(envSecretKey)
	if bucketName == "" || accessKeyID == "" || secretAccessKey == "" {
		fmt.Println("Miss environment variables for SourceQiniu test")
		t.SkipNow()
	}
	source, err := NewSourceQiniu(bucketName, accessKeyID, secretAccessKey)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	logger := utils.GetLogger()
	// Clear test bucket
	mac := &digest.Mac{AccessKey: accessKeyID, SecretKey: []byte(secretAccessKey)}
	qnClient := rs.New(mac)
	clearBucket(&qnClient, bucketName)

	for _, test := range sourceQiniuTests {
		for _, object := range test.SourceObjects {
			putPolicy := rs.PutPolicy{
				Scope: bucketName,
			}
			token := putPolicy.Token(nil)
			str := "hello world"
			err := qnio.Put2(nil, nil, token, object, bytes.NewBuffer([]byte(str)), int64(len(str)), nil)
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

		clearBucket(&qnClient, bucketName)
		clearSourceQiniu(source)
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

func clearBucket(client *rs.Client, bucketName string) {
	service := rsf.New(nil)
	objects, _, _ := service.ListPrefix(nil, bucketName, "", "", -1)
	for _, object := range objects {
		client.Delete(nil, bucketName, object.Key)
	}
}

func clearSourceQiniu(source *SourceQiniu) {
	source.lastModifiedTimes = make(map[string]time.Time)
	source.marker = ""
}
