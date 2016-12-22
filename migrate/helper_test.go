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

package migrate

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/yunify/qingstor-sdk-go/config"
	"github.com/yunify/qingstor-sdk-go/service"
)

const (
	testSrcListPath = "/tmp/test_src_list"
	testConfigPath  = "/tmp/test_qs_config"
	testRecordPath  = "/tmp/test_record"
	testQSAccessKey = "test_qs_access_key"
	testQSSecretKey = "test_qs_secret_key"

	publicReadACL = `
{
  "acl": [
    {
      "grantee": {
          "type": "group",
          "name": "QS_ALL_USERS"
      },
      "permission": "READ"
    }
  ]
}`
	tmpMoveSuffix = ".move"
)

var testConfigContent = fmt.Sprintf(
	"access_key_id: %s\nsecret_access_key: %s\n",
	testQSAccessKey, testQSSecretKey,
)

// Remember to call tearDownTestFiles when finishing tests.
func setUpTestFiles(t *testing.T, configContent string) {
	if configContent == "" {
		return
	}
	os.Create(testSrcListPath)

	configFile, err := os.Create(testConfigPath)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	configFile.WriteString(configContent)

	defaultConfPath := config.GetUserConfigFilePath()
	err = os.MkdirAll(path.Dir(defaultConfPath), 0644)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	// If default config file exists, move it temporarily.
	if _, err := os.Stat(defaultConfPath); err == nil {
		tmpMovePath := defaultConfPath + tmpMoveSuffix
		os.Rename(defaultConfPath, tmpMovePath)
	}
	defaultConfFile, err := os.Create(defaultConfPath)
	defaultConfFile.Write([]byte(configContent))
}

func tearDownTestFiles() {
	os.Remove(testSrcListPath)
	os.Remove(testRecordPath)
	os.Remove(testConfigPath)
	defaultConfPath := config.GetUserConfigFilePath()
	os.Remove(defaultConfPath)

	tmpMovePath := defaultConfPath + tmpMoveSuffix
	if _, err := os.Stat(tmpMovePath); err == nil {
		os.Rename(tmpMovePath, defaultConfPath)
	}
}

func clearBucket(bucket *service.Bucket) {
	output, err := bucket.ListObjects(&service.ListObjectsInput{})
	if err != nil {
		fmt.Println(err)
	} else {
		for _, object := range output.Keys {
			bucket.DeleteObject(object.Key)
		}
	}
	_, err = bucket.Delete()
	if err != nil {
		fmt.Println(err)
	}
}

func setBucketACLPublicRead(bucket *service.Bucket) error {
	putBucketACLInput := &service.PutBucketACLInput{}
	err := json.Unmarshal([]byte(publicReadACL), putBucketACLInput)
	if err != nil {
		return err
	}
	bucket.PutACL(putBucketACLInput)
	return nil
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
