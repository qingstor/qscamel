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

package record

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Recorder(t *testing.T) {
	recordA := "http://example.com/A"
	recordB := "http://example.com/B"
	recordC := "http://example.com/C"
	existingRecords := []string{recordA, recordB}

	testFile := "/tmp/test_record"
	file, err := os.Create(testFile)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	for _, record := range existingRecords {
		file.WriteString(fmt.Sprintf("%s\n", record))
	}
	file.Seek(0, 0)

	r := NewRecorder(file)
	assert.Equal(t, existingRecords, []string(r.records))
	assert.True(t, r.IsExist(recordA))
	assert.True(t, r.IsExist(recordB))
	assert.False(t, r.IsExist(recordC))
	r.Put(recordC)
	file.Close()

	// Create Record struct again to see if recordC has been written.
	file, err = os.Open(testFile)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	r = NewRecorder(file)
	assert.Equal(
		t, append(existingRecords, recordC), []string(r.records),
	)
	os.Remove(testFile)
}
