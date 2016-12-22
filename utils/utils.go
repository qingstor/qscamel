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

package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	log "github.com/frostyplanet/logrus"

	"github.com/yunify/qingstor-sdk-go/service"
)

// ExpandHomeDirectory expands tilde in relatedPath to concrete home path.
func ExpandHomeDirectory(relatedPath string) string {
	var home string
	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = os.Getenv("HOME")
	}

	return strings.Replace(relatedPath, "~", home, 1)
}

// GetLogger creates logger using TextFormatter.
func GetLogger() *log.Logger {
	logger := log.New()
	logger.Formatter = &log.TextFormatter{
		FullTimestamp: true, TimestampFormat: time.RFC822,
	}
	return logger
}

func printDryRun(dryRun bool) string {
	if dryRun {
		return " (DRY RUN)"
	}
	return ""
}

// LogResult writes processed source sites to logger.
// The type of source sites depends on situation.
func LogResult(sourceSites []string, situation string, logger *log.Logger, dryRun bool) {
	if len(sourceSites) > 0 {
		logger.Infof("%s migration:%s", strings.Title(situation), printDryRun(dryRun))
		for _, sourceSite := range sourceSites {
			logger.Infoln(sourceSite)
		}
	} else {
		logger.Infof("No %s migration.%s", situation, printDryRun(dryRun))
	}
}

// DetermineBucketZone determines the zone of bucket using List Buckets api.
func DetermineBucketZone(qsService *service.Service, bucketName string) (string, error) {
	output, err := qsService.ListBuckets(nil)
	if err != nil {
		return "", err
	}
	for _, bucket := range output.Buckets {
		if bucketName == bucket.Name {
			return bucket.Location, nil
		}
	}
	return "", fmt.Errorf("cannot determine which zone the bucket \"%s\" resides", bucketName)
}

// CheckErrorForExit checks if error occurs.
// If error is not nil, print the error message and exit the application.
// If error is nil, do nothing.
func CheckErrorForExit(err error, code ...int) {
	if err != nil {
		exitCode := 1
		if len(code) > 0 {
			exitCode = code[0]
		}
		fmt.Println(err.Error(), exitCode)
		os.Exit(exitCode)
	}
}
