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
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/upyun/go-sdk/upyun"
	"github.com/yunify/qscamel/record"
	"github.com/yunify/qscamel/utils"
)

const (
	envBucketName = "QSCAMEL_TEST_UPYUN_BUCKETNAME"
	envZone       = "QSCAMEL_TEST_UPYUN_ZONE"
	envAccessKey  = "QSCAMEL_TEST_UPYUN_ACCESS_KEY_ID"
	envSecretKey  = "QSCAMEL_TEST_UPYUN_SECRET_ACCESS_KEY"
)

type testCase struct {
	SourceObjects       []string
	ThreadNum           int
	RecordObjects       []string
	ExpectedObjectNames []string
	ExpectedSkipped     []string
	ExpectedDone        bool
}

var sourceUpyunTests = []testCase{
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
		SourceObjects:       []string{"object4", "object5", "object6"},
		ThreadNum:           5,
		ExpectedObjectNames: []string{"object4", "object5", "object6"},
		ExpectedDone:        true,
	},

	// Skip recorded objects
	{
		SourceObjects:       []string{"object7", "object8", "object9"},
		ThreadNum:           3,
		RecordObjects:       []string{"object7"},
		ExpectedObjectNames: []string{"object8", "object9"},
		ExpectedSkipped:     []string{"object7"},
		ExpectedDone:        true,
	},
}

func TestSourceUpyun_GetSourceSites(t *testing.T) {
	bucketName := os.Getenv(envBucketName)
	zone := os.Getenv(envZone)
	accessKey := os.Getenv(envAccessKey)
	secretKey := os.Getenv(envSecretKey)
	if zone == "" || accessKey == "" || secretKey == "" {
		fmt.Println("Miss environment variables for NewSourceUpyun test")
		t.SkipNow()
	}

	source, err := NewSourceUpyun(bucketName, zone, accessKey, secretKey)
	assert.Nil(t, err)

	logger := utils.GetLogger()
	for _, test := range sourceUpyunTests {
		for _, object := range test.SourceObjects {

			s := "hello world"
			r := strings.NewReader(s)

			err := source.up.Put(&upyun.PutObjectConfig{
				Path:   "/" + object,
				Reader: r,
				Headers: map[string]string{
					"Content-Length": fmt.Sprint(len(s)),
				},
			})
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

		sort.Strings(sourceSites)
		sort.Strings(names)

		assert.Nil(t, err)
		expectedSourceSites := getSourceSites(source.BucketURL, test.ExpectedObjectNames)
		assert.Equal(t, expectedSourceSites, sourceSites)
		assert.Equal(t, test.ExpectedObjectNames, names)
		expectedSkipped := getSourceSites(source.BucketURL, test.ExpectedSkipped)
		assert.Equal(t, expectedSkipped, skipped)
		assert.Equal(t, test.ExpectedDone, done)

		clearBucket(source.up)
		clearSourceUpyun(source)
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

func clearBucket(up *upyun.UpYun) {

	objsChan := make(chan *upyun.FileInfo, 10)
	err := up.List(&upyun.GetObjectsConfig{
		Path:         "/",
		ObjectsChan:  objsChan,
		MaxListLevel: -1,
	})

	if err == nil {
		for object := range objsChan {
			err = up.Delete(&upyun.DeleteObjectConfig{
				Path: object.Name,
			})
		}
	}

}

func clearSourceUpyun(source *SourceUpyun) {
	source.lastModifiedTimes = make(map[string]time.Time)
	source.objsChan = nil
	source.marker = ""
}
