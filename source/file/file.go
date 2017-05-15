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
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/frostyplanet/logrus"

	"github.com/yunify/qscamel/record"
)

const objectNameSep = " "

// SourceFile is source-list file type source
type SourceFile struct {
	SourceList *os.File
	client     *http.Client
}

// NewSourceFile creates an instance of SourceFile
func NewSourceFile(listPath string) (*SourceFile, error) {
	sourceListFile, err := os.Open(listPath)
	if err != nil {
		return &SourceFile{}, fmt.Errorf("can't open source list %s (%s)", listPath, err.Error())
	}

	return &SourceFile{
		SourceList: sourceListFile,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Second * 10,
		},
	}, nil
}

// GetSourceSites implements MigrateSource.GetSourceSites
func (source *SourceFile) GetSourceSites(
	threadNum int, logger *log.Logger, recorder *record.Recorder,
) (sourceSites []string, objectNames []string, skipped []string, done bool, err error) {
	sourceSites, objectNames, skipped = []string{}, []string{}, []string{}
	reader := bufio.NewReader(source.SourceList)
	var objectName string
	for i := 0; i < threadNum; {
		sourceSite, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
				done = true
			}
			return sourceSites, objectNames, skipped, done, err
		}
		if sourceSite == "\n" || sourceSite[0] == '#' {
			continue
		}
		sourceSite = sourceSite[:len(sourceSite)-1]
		if strings.Contains(sourceSite, objectNameSep) {
			// Format: source site [spacing] object name
			results := strings.Split(sourceSite, objectNameSep)
			if len(results) == 2 {
				sourceSite = results[0]
				objectName = results[1]
			} else {
				logger.Warnf(
					"Skip invalid source file format %s",
					sourceSite,
				)
				skipped = append(skipped, sourceSite)
				continue
			}
		} else {
			// Format: just source site
			// Use source site's path as object name.
			sourceURL, _ := url.Parse(sourceSite)
			if strings.HasPrefix(sourceURL.Path, "/") {
				objectName = sourceURL.Path[1:]
			} else {
				objectName = sourceURL.Host
				logger.Warnf(
					"Source site %s has no path, "+
						"use host as object name%s",
					sourceSite,
					sourceURL.Host,
				)
			}
		}
		//TODO check both source site and object name
		if recorder.IsExist(sourceSite) {
			logger.Infof(
				"Skip completed source site: %s", sourceSite,
			)
			skipped = append(skipped, sourceSite)
			continue
		}
		sourceSites = append(sourceSites, sourceSite)
		objectNames = append(objectNames, objectName)
		i++
	}
	return
}

// GetSourceSiteInfo implements MigrateSource.GetSourceSites
func (source *SourceFile) GetSourceSiteInfo(sourceSite string) (lastModified time.Time, err error) {
	resp, err := source.client.Get(sourceSite)
	if resp != nil && !resp.Close {
		resp.Body.Close()
	}
	if err != nil {
		return time.Time{}, err
	}
	date := resp.Header.Get("Last-Modified")
	lastModified, err = time.Parse(time.RFC1123, date)
	return
}
