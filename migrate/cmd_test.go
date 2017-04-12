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
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qingstor-sdk-go/config"
	"github.com/yunify/qscamel/utils"
)

var initContextTests = []struct {
	Flags         []string
	ConfigContent string
	IfPassCheck   bool
}{
	// Expect success, full name flags
	{
		[]string{
			"--bucket=test-qs-bucket",
			`--description="test description"`,
			"--src-type=file",
			fmt.Sprintf("--src=%s", testSrcListPath),
			fmt.Sprintf("--config=%s", testConfigPath),
		},
		testConfigContent, true,
	},
	// Expect success, shorthand flags
	{
		[]string{
			"-b", "test-qs-bucket",
			"-d", "test_description",
			"-t", "file",
			"-s", testSrcListPath,
			"-c", testConfigPath,
		},
		testConfigContent, true,
	},
	// Expect success, use default config file
	{
		[]string{
			"--bucket=test-qs-bucket",
			`--description="test description"`,
			"--src-type=file",
			fmt.Sprintf("--src=%s", testSrcListPath),
		},
		testConfigContent, true,
	},
	// Expect success, print version
	{
		[]string{
			"--version",
		},
		"", true,
	},
	// Expect failure, miss QingStor bucket
	{
		[]string{
			`--description="test description"`,
			"--src-type=file",
			fmt.Sprintf("--src=%s", testSrcListPath),
			fmt.Sprintf("--config=%s", testConfigPath),
		},
		testConfigContent, false,
	},
	// Expect failure, miss description
	{
		[]string{
			"--bucket=test-qs-bucket",
			"--src-type=file",
			fmt.Sprintf("--src=%s", testSrcListPath),
			fmt.Sprintf("--config=%s", testConfigPath),
		},
		testConfigContent, false,
	},
	// Expect failure, miss source type
	{
		[]string{
			"--bucket=test-qs-bucket",
			fmt.Sprintf("--src=%s", testSrcListPath),
			`--description="test description"`,
			fmt.Sprintf("--config=%s", testConfigPath),
		},
		testConfigContent, false,
	},
	// Expect failure, miss specific source
	{
		[]string{
			"--bucket=test-qs-bucket",
			"--src-type=file",
			`--description="test description"`,
			fmt.Sprintf("--config=%s", testConfigPath),
		},
		testConfigContent, false,
	},
	// Expect failure, source list file doesn't exist
	{
		[]string{
			"--bucket=test-qs-bucket",
			`--description="test description"`,
			"--src-type=file",
			"--src=/tmp/non_existing_test_src_list",
			fmt.Sprintf("--config=%s", testConfigPath),
		},
		testConfigContent, false,
	},
	// Expect failure, config file doesn't exist
	{
		[]string{
			"--bucket=test-qs-bucket",
			`--description="test description"`,
			"--src-type=file",
			fmt.Sprintf("--src=%s", testSrcListPath),
			"--config=/tmp/non_existing_test_qs_config",
		},
		testConfigContent, false,
	},
	// Expect failure, incomplete yaml configuration
	{
		[]string{
			"--bucket=test-qs-bucket",
			`--description="test description"`,
			"--src-type=file",
			fmt.Sprintf("--src=%s", testSrcListPath),
			fmt.Sprintf("--config=%s", testConfigPath),
		},
		fmt.Sprintf(
			"access_key_id: %s\n", testQSAccessKey,
		), false,
	},
}

func resetFlagVariables() {
	ctx = &Context{
		IgnoreUnmodified: true,
		Logger:           utils.GetLogger(),
		QSConfig:         &config.Config{Connection: &http.Client{}},
	}
	specificSource = ""
	sourceZone = ""
	sourceAccessKeyID = ""
	sourceSecretAccessKey = ""
	configPath = ""
	logPath = ""
	description = ""
}

func Test_CheckFlags(t *testing.T) {
	program := os.Args[0]

	for _, test := range initContextTests {
		setUpTestFiles(t, test.ConfigContent)
		// Reset flag variables, command-line arguments and cobra.Command.
		resetFlagVariables()
		os.Args = []string{program}
		os.Args = append(os.Args, test.Flags...)
		cmd := &cobra.Command{
			PreRunE: checkFlags,
			Run:     func(cmd *cobra.Command, args []string) {},
		}
		defineFlags(cmd)

		err := cmd.Execute()
		assert.Equal(t, test.IfPassCheck, err == nil)
	}
	tearDownTestFiles()
}

// Set all flags and check fields of Context.
func Test_CheckFlags_AssertContext(t *testing.T) {
	sourceType := "file"
	qsBucketName := "test-qs-bucket"
	threadNum := 15
	logPath := "/tmp/log_file"
	flags := []string{
		"-b", qsBucketName,
		"-s", testSrcListPath,
		"-t", sourceType,
		"-c", testConfigPath,
		"-T", strconv.Itoa(threadNum),
		"-l", logPath,
		"-d", "migrate 10",
		"-n",
	}

	setUpTestFiles(t, testConfigContent)
	resetFlagVariables()
	program := os.Args[0]
	os.Args = []string{program}
	os.Args = append(os.Args, flags...)
	cmd := &cobra.Command{
		PreRunE: checkFlags,
		Run:     func(cmd *cobra.Command, args []string) {},
	}
	defineFlags(cmd)
	err := cmd.Execute()
	assert.Nil(t, err)
	assert.Equal(t, ctx.SourceType, sourceType)
	assert.Equal(t, testQSAccessKey, ctx.QSConfig.AccessKeyID)
	assert.Equal(t, testQSSecretKey, ctx.QSConfig.SecretAccessKey)
	assert.Equal(t, qsBucketName, ctx.QSBucketName)
	assert.Equal(t, threadNum, ctx.ThreadNum)
	assert.NotNil(t, ctx.Logger.Out)
	assert.NotNil(t, ctx.Recorder)
	assert.True(t, ctx.IgnoreUnmodified)
	assert.True(t, ctx.DryRun)
	tearDownTestFiles()
	if ctx.Recorder != nil {
		ctx.Recorder.Clear()
	}
}

func Test_CheckFlags_ThreadNumOverLimit(t *testing.T) {
	sourceType := "file"
	qsBucketName := "test-qs-bucket"
	threadNum := MaxThreadNum + 50
	flags := []string{
		"-b", qsBucketName,
		"-s", testSrcListPath,
		"-t", sourceType,
		"-d", "thread num over limit",
		"-T", strconv.Itoa(threadNum),
	}

	setUpTestFiles(t, testConfigContent)
	resetFlagVariables()
	program := os.Args[0]
	os.Args = []string{program}
	os.Args = append(os.Args, flags...)
	cmd := &cobra.Command{
		PreRunE: checkFlags,
		Run:     func(cmd *cobra.Command, args []string) {},
	}
	defineFlags(cmd)
	err := cmd.Execute()
	assert.Nil(t, err)
	assert.Equal(t, ctx.ThreadNum, MaxThreadNum)
	tearDownTestFiles()
	if ctx.Recorder != nil {
		ctx.Recorder.Clear()
	}
}

func Test_CheckFlags_ThreadNumSetToDefault(t *testing.T) {
	sourceType := "file"
	qsBucketName := "test-qs-bucket"
	flags := []string{
		"-b", qsBucketName,
		"-s", testSrcListPath,
		"-t", sourceType,
		"-d", "thread num not set, use default instead",
	}

	setUpTestFiles(t, testConfigContent)
	resetFlagVariables()
	program := os.Args[0]
	os.Args = []string{program}
	os.Args = append(os.Args, flags...)
	cmd := &cobra.Command{
		PreRunE: checkFlags,
		Run:     func(cmd *cobra.Command, args []string) {},
	}
	defineFlags(cmd)
	err := cmd.Execute()
	assert.Nil(t, err)
	assert.Equal(t, ctx.ThreadNum, DefaultThreadNum)
	tearDownTestFiles()
	if ctx.Recorder != nil {
		ctx.Recorder.Clear()
	}
}
