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
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qingstor-sdk-go/config"
	"github.com/yunify/qingstor-sdk-go/service"
	"github.com/yunify/qscamel/record"
	"github.com/yunify/qscamel/source/file"
	"github.com/yunify/qscamel/utils"
)

const (
	envAccessKey   = "QSCAMEL_TEST_ACCESS_KEY_ID"
	envSecretKey   = "QSCAMEL_TEST_SECRET_ACCESS_KEY"
	testZone       = "pek3a"
	testSourceType = "file"
)

// Read environment variables, write to default config file and init Config.
func initQSConfig(t *testing.T) *config.Config {
	accessKey := os.Getenv(envAccessKey)
	secretKey := os.Getenv(envSecretKey)
	if accessKey == "" || secretKey == "" {
		fmt.Printf(
			"Miss environment variables %s or %s for test",
			envAccessKey, envSecretKey,
		)
		t.SkipNow()
	}
	configContent := fmt.Sprintf(
		"access_key_id: %s\nsecret_access_key: %s\n",
		accessKey, secretKey,
	)
	setUpTestFiles(t, configContent)
	c := &config.Config{Connection: &http.Client{}}
	err := c.LoadConfigFromFilepath(config.GetUserConfigFilePath())
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	return c
}

func setUpBucket(t *testing.T, c *config.Config) *service.Bucket {
	bucketName := "test-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	qsService, _ := service.Init(c)
	bucket, _ := qsService.Bucket(bucketName, testZone)
	_, err := bucket.Put()
	assert.Nil(t, err)
	return bucket
}

// Use a public-readable bucket as source site.
func setUpUpstreamBucket(t *testing.T, c *config.Config) (*service.Bucket, string) {
	bucketURL := "http://"
	bucket := setUpBucket(t, c)
	err := setBucketACLPublicRead(bucket)
	assert.Nil(t, err)
	bucketURL += strings.Join(
		[]string{bucket.Properties.BucketName, testZone, c.Host}, ".",
	)
	return bucket, bucketURL
}

// Put sample objects.
func putObjects(t *testing.T, bucket *service.Bucket, objects []string) {
	content := "hello world"
	for _, name := range objects {
		_, err := bucket.PutObject(
			name, &service.PutObjectInput{
				ContentLength: len(content),
				Body:          bytes.NewBuffer([]byte(content)),
			})
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
	}
}

// Write source sites for source list file and record file.
func setUpSourceSites(t *testing.T, fileName, upstreamURL string, objects []string) {
	file, err := os.Create(fileName)
	assert.Nil(t, err)
	for _, object := range objects {
		file.WriteString(
			fmt.Sprintf("%s/%s\n", upstreamURL, object),
		)
	}
	file.Close()
}

func initCtx(t *testing.T, source string, c *config.Config, bucketName, recordPath string) *Context {
	s, err := file.NewSourceFile(source)
	assert.Nil(t, err)
	recordFile, err := os.Open(recordPath)
	assert.Nil(t, err)
	return &Context{
		SourceType:   testSourceType,
		Source:       s,
		QSConfig:     c,
		QSBucketName: bucketName,
		// default
		IgnoreUnmodified: true,
		ThreadNum:        DefaultThreadNum,
		Logger:           utils.GetLogger(),
		Recorder:         record.NewRecorder(recordFile),
	}
}

// Complete normal source sites, fail in non-existing one and skip recorded source sites.
func Test_Migrate_SkipRecord(t *testing.T) {
	sourceObjects := []string{"object1", "object2", "object3", "object4"}
	recordedObjects := []string{"object1", "object2"}
	nonExistingObjects := []string{"object5"}
	expectedCompleted := []string{"object3", "object4"}
	expectedFailed := nonExistingObjects
	expectedSkipped := recordedObjects

	c := initQSConfig(t)
	upstream, upstreamURL := setUpUpstreamBucket(t, c)
	putObjects(t, upstream, sourceObjects)
	bucket := setUpBucket(t, c)
	setUpSourceSites(
		t, testSrcListPath, upstreamURL, append(sourceObjects, nonExistingObjects...),
	)
	setUpSourceSites(t, testRecordPath, upstreamURL, recordedObjects)
	context := initCtx(t, testSrcListPath, c, bucket.Properties.BucketName, testRecordPath)

	completed, failed, skipped, err := Migrate(context)
	assert.Nil(t, err)

	expectedCompleted = getSourceSites(upstreamURL, expectedCompleted)
	expectedFailed = getSourceSites(upstreamURL, expectedFailed)
	expectedSkipped = getSourceSites(upstreamURL, expectedSkipped)
	for _, s := range [][]string{completed, failed, skipped,
		expectedCompleted, expectedFailed, expectedSkipped} {
		sort.Strings(s)
	}
	assert.Equal(t, expectedCompleted, completed)
	assert.Equal(t, expectedFailed, failed)
	assert.Equal(t, expectedSkipped, skipped)

	clearBucket(upstream)
	clearBucket(bucket)
	tearDownTestFiles()
}

// Test RryRun: completed objects hasn't been actually migrated.
func Test_Migrate_DryRun(t *testing.T) {
	sourceObjects := []string{"object1", "object2", "object3", "object4"}
	recordedObjects := []string{"object1", "object2"}
	expectedCompleted := []string{"object3", "object4"}
	expectedFailed := []string{}
	expectedSkipped := recordedObjects

	c := initQSConfig(t)
	upstream, upstreamURL := setUpUpstreamBucket(t, c)
	putObjects(t, upstream, sourceObjects)
	bucket := setUpBucket(t, c)
	setUpSourceSites(t, testSrcListPath, upstreamURL, sourceObjects)
	setUpSourceSites(t, testRecordPath, upstreamURL, recordedObjects)
	context := initCtx(t, testSrcListPath, c, bucket.Properties.BucketName, testRecordPath)
	context.DryRun = true

	completed, failed, skipped, err := Migrate(context)
	assert.Nil(t, err)

	expectedCompleted = getSourceSites(upstreamURL, expectedCompleted)
	expectedFailed = getSourceSites(upstreamURL, expectedFailed)
	expectedSkipped = getSourceSites(upstreamURL, expectedSkipped)
	for _, s := range [][]string{completed, failed, skipped,
		expectedCompleted, expectedFailed, expectedSkipped} {
		sort.Strings(s)
	}
	assert.Equal(t, expectedCompleted, completed)
	assert.Equal(t, expectedFailed, failed)
	assert.Equal(t, expectedSkipped, skipped)

	// Expect bucket has no objects.
	output, err := bucket.ListObjects(&service.ListObjectsInput{})
	assert.Equal(t, 0, len(output.Keys))

	clearBucket(upstream)
	clearBucket(bucket)
	tearDownTestFiles()
}

// If context.IgnoreUnmodified is true(default), ignore latest objects and
// fetch modified objects.
func Test_Migrate_IgnoreUnmodified(t *testing.T) {
	sourceObjects := []string{"object1", "object2", "object3", "object4"}
	priorObjects := []string{"object1"}
	latestObjects := []string{"object2"}
	expectedCompleted := []string{"object1", "object3", "object4"}
	expectedFailed := []string{}
	expectedSkipped := latestObjects

	c := initQSConfig(t)
	// Modified time of object is prior to upstream
	bucket := setUpBucket(t, c)
	putObjects(t, bucket, priorObjects)
	time.Sleep(time.Second)
	// Upload objects to upstream
	upstream, upstreamURL := setUpUpstreamBucket(t, c)
	putObjects(t, upstream, sourceObjects)
	time.Sleep(time.Second)
	// Modified time of object is after upstream
	putObjects(t, bucket, latestObjects)
	setUpSourceSites(t, testSrcListPath, upstreamURL, sourceObjects)
	setUpSourceSites(t, testRecordPath, upstreamURL, []string{})
	context := initCtx(t, testSrcListPath, c, bucket.Properties.BucketName, testRecordPath)

	completed, failed, skipped, err := Migrate(context)
	assert.Nil(t, err)

	expectedCompleted = getSourceSites(upstreamURL, expectedCompleted)
	expectedFailed = getSourceSites(upstreamURL, expectedFailed)
	expectedSkipped = getSourceSites(upstreamURL, expectedSkipped)
	for _, s := range [][]string{completed, failed, skipped,
		expectedCompleted, expectedFailed, expectedSkipped} {
		sort.Strings(s)
	}
	assert.Equal(t, expectedCompleted, completed)
	assert.Equal(t, expectedFailed, failed)
	assert.Equal(t, expectedSkipped, skipped)

	clearBucket(upstream)
	clearBucket(bucket)
	tearDownTestFiles()
}

// If context.Overwrite is true, overwrite existing objects.
func Test_Migrate_Overwrite(t *testing.T) {
	sourceObjects := []string{"object1", "object2", "object3", "object4"}
	existingObjects := []string{"object1", "object2"}
	expectedCompleted := sourceObjects
	expectedFailed := []string{}
	expectedSkipped := []string{}

	c := initQSConfig(t)
	upstream, upstreamURL := setUpUpstreamBucket(t, c)
	putObjects(t, upstream, sourceObjects)
	bucket := setUpBucket(t, c)
	putObjects(t, bucket, existingObjects)
	setUpSourceSites(t, testSrcListPath, upstreamURL, sourceObjects)
	setUpSourceSites(t, testRecordPath, upstreamURL, []string{})
	context := initCtx(t, testSrcListPath, c, bucket.Properties.BucketName, testRecordPath)
	context.Overwrite = true

	completed, failed, skipped, err := Migrate(context)
	assert.Nil(t, err)

	expectedCompleted = getSourceSites(upstreamURL, expectedCompleted)
	expectedFailed = getSourceSites(upstreamURL, expectedFailed)
	expectedSkipped = getSourceSites(upstreamURL, expectedSkipped)
	for _, s := range [][]string{completed, failed, skipped,
		expectedCompleted, expectedFailed, expectedSkipped} {
		sort.Strings(s)
	}
	assert.Equal(t, expectedCompleted, completed)
	assert.Equal(t, expectedFailed, failed)
	assert.Equal(t, expectedSkipped, skipped)

	clearBucket(upstream)
	clearBucket(bucket)
	tearDownTestFiles()
}

// If context.IgnoreExisting is true, ignore existing objects.
func Test_Migrate_IgnoreExisting(t *testing.T) {
	sourceObjects := []string{"object1", "object2", "object3", "object4"}
	existingObjects := []string{"object1", "object2"}
	expectedCompleted := []string{"object3", "object4"}
	expectedFailed := []string{}
	expectedSkipped := existingObjects

	c := initQSConfig(t)
	upstream, upstreamURL := setUpUpstreamBucket(t, c)
	putObjects(t, upstream, sourceObjects)
	bucket := setUpBucket(t, c)
	putObjects(t, bucket, existingObjects)
	setUpSourceSites(t, testSrcListPath, upstreamURL, sourceObjects)
	setUpSourceSites(t, testRecordPath, upstreamURL, []string{})
	context := initCtx(t, testSrcListPath, c, bucket.Properties.BucketName, testRecordPath)
	context.IgnoreExisting = true

	completed, failed, skipped, err := Migrate(context)
	assert.Nil(t, err)

	expectedCompleted = getSourceSites(upstreamURL, expectedCompleted)
	expectedFailed = getSourceSites(upstreamURL, expectedFailed)
	expectedSkipped = getSourceSites(upstreamURL, expectedSkipped)
	for _, s := range [][]string{completed, failed, skipped,
		expectedCompleted, expectedFailed, expectedSkipped} {
		sort.Strings(s)
	}
	assert.Equal(t, expectedCompleted, completed)
	assert.Equal(t, expectedFailed, failed)
	assert.Equal(t, expectedSkipped, skipped)

	clearBucket(upstream)
	clearBucket(bucket)
	tearDownTestFiles()
}
