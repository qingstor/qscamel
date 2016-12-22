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
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

// Recorder holds string record via slice and writes record to file for persistence.
type Recorder struct {
	file    *os.File
	records sort.StringSlice
}

// NewRecorder reads records from file and sort records.
func NewRecorder(file *os.File) *Recorder {
	r := &Recorder{file: file, records: []string{}}
	reader := bufio.NewReader(file)
	for {
		record, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		r.records = append(r.records, record[:len(record)-1])
	}
	r.records.Sort()
	return r
}

// Put appends record to slice and writes it to file.
func (r *Recorder) Put(record string) {
	r.records = append(r.records, record)
	r.file.WriteString(fmt.Sprintf("%s\n", record))
}

// IsExist checks if a record exists in slice.
func (r *Recorder) IsExist(record string) bool {
	pos := r.records.Search(record)
	return pos < len(r.records) && r.records[pos] == record
}

// Clear removes record file and record string slice.
func (r *Recorder) Clear() {
	os.Remove(r.file.Name())
	r.records = []string{}
}
