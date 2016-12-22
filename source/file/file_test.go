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

package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yunify/qscamel/record"
	"github.com/yunify/qscamel/utils"
)

var getSourceSitesTests = []struct {
	SourceSites         []string
	ThreadNum           int
	CompletedRecords    []string
	ExpectedSourceSites []string
	ExpectedObjectNames []string
	ExpectedSkipped     []string
	ExpectedDone        bool
}{
	// Format: just source site
	{
		[]string{
			"https://image.example.com/public/cat.jpg\n",
			"https://image.example.com/public/dog.jpg\n",
		},
		2, nil,
		[]string{
			"https://image.example.com/public/cat.jpg",
			"https://image.example.com/public/dog.jpg",
		},
		[]string{"public/cat.jpg", "public/dog.jpg"},
		[]string{}, false,
	},
	// Format: source site [spacing] object name
	{
		[]string{
			"https://image.example.com/public/cat.jpg image1.jpg\n",
			"https://image.example.com/public/dog.jpg image2.jpg\n",
		},
		2, nil,
		[]string{
			"https://image.example.com/public/cat.jpg",
			"https://image.example.com/public/dog.jpg",
		},
		[]string{"image1.jpg", "image2.jpg"},
		[]string{}, false,
	},
	// EOF before reaching threadNum
	{
		[]string{
			"https://image.example.com/public/cat.jpg\n",
			"https://image.example.com/public/dog.jpg\n",
		},
		4, nil,
		[]string{
			"https://image.example.com/public/cat.jpg",
			"https://image.example.com/public/dog.jpg",
		},
		[]string{"public/cat.jpg", "public/dog.jpg"},
		[]string{}, true,
	},
	// Source site has no path, use host.
	{
		[]string{
			"https://image.example.com\n",
			"https://image.example.com/public/dog.jpg\n",
		},
		2, nil,
		[]string{
			"https://image.example.com",
			"https://image.example.com/public/dog.jpg",
		},
		[]string{"image.example.com", "public/dog.jpg"},
		[]string{}, false,
	},
	// Skip invalid line format
	{
		[]string{
			"https://image.example.com/public/cat.jpg image1.jpg image.jpg\n",
			"https://image.example.com/public/dog.jpg image2.jpg\n",
		},
		2, nil,
		[]string{
			"https://image.example.com/public/dog.jpg",
		},
		[]string{"image2.jpg"},
		[]string{
			"https://image.example.com/public/cat.jpg image1.jpg image.jpg",
		},
		true,
	},
	// Skip completed source site from recorder
	{
		[]string{
			"https://image.example.com/public/cat.jpg image1.jpg\n",
			"https://image.example.com/public/dog.jpg image2.jpg\n",
		},
		2,
		// CompletedRecords
		[]string{"https://image.example.com/public/dog.jpg"},
		[]string{"https://image.example.com/public/cat.jpg"},
		[]string{"image1.jpg"},
		[]string{
			"https://image.example.com/public/dog.jpg",
		}, true,
	},
}

func Test_GetSourceSites(t *testing.T) {
	testList := "/tmp/migration_source_list"
	logger := utils.GetLogger()
	for _, test := range getSourceSitesTests {
		// Write source list file.
		file, err := os.Create(testList)
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		for _, sourceSite := range test.SourceSites {
			file.WriteString(sourceSite)
		}
		file.Close()

		// Init SourceFile and Recorder.
		source, err := NewSourceFile(testList)
		assert.Nil(t, err)
		recordFile, err := record.GetRecordFile("test_get_source_sites")
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		if test.CompletedRecords != nil {
			for _, record := range test.CompletedRecords {
				recordFile.WriteString(fmt.Sprintf("%s\n", record))
			}
			recordFile.Seek(0, 0)
		}
		recorder := record.NewRecorder(recordFile)

		sourceSites, names, skipped, done, err := source.GetSourceSites(
			test.ThreadNum, logger, recorder,
		)
		assert.Nil(t, err)
		assert.Equal(t, test.ExpectedSourceSites, sourceSites)
		assert.Equal(t, test.ExpectedObjectNames, names)
		assert.Equal(t, test.ExpectedSkipped, skipped)
		assert.Equal(t, test.ExpectedDone, done)
		os.Remove(recordFile.Name())
	}
	os.Remove(testList)
}

// Test skipping black lines and comment lines
func Test_GetSourceSites_Skip(t *testing.T) {
	testList := "/tmp/migration_source_list"
	logger := utils.GetLogger()
	sourceSiteA := "https://image.example.com/public/cat.jpg"
	sourceSiteB := "https://image.example.com/public/dog.jpg"
	threadNum := 2

	file, err := os.Create(testList)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	/*Test source list content:
	  https://image.example.com/public/cat.jpg


	  # This is a comment
	  https://image.example.com/public/dog.jpg
	*/
	file.WriteString(fmt.Sprintf("%s\n", sourceSiteA))
	file.Write([]byte("\n\n"))
	file.Write([]byte("# This is a comment\n"))
	file.WriteString(fmt.Sprintf("%s\n", sourceSiteB))
	file.Close()
	source, err := NewSourceFile(testList)
	assert.Nil(t, err)
	recordFile, err := record.GetRecordFile("test_get_source_sites_skip")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	recorder := record.NewRecorder(recordFile)

	sourceSites, _, failed, done, err := source.GetSourceSites(threadNum, logger, recorder)
	assert.Nil(t, err)
	assert.Equal(t, []string{sourceSiteA, sourceSiteB}, sourceSites)
	assert.Equal(t, []string{}, failed)
	assert.Equal(t, false, done)
	os.Remove(testList)
	os.Remove(recordFile.Name())
}

func Test_GetSourceSiteInfo(t *testing.T) {
	testList := "/tmp/migration_source_list"
	file, err := os.Create(testList)
	assert.Nil(t, err)
	file.Close()
	source, err := NewSourceFile(testList)
	assert.Nil(t, err)

	// Can't connect to source site, expect error.
	sourceSite := "http://notexsitingsourcesite.com"
	_, err = source.GetSourceSiteInfo(sourceSite)
	assert.NotNil(t, err)
}
