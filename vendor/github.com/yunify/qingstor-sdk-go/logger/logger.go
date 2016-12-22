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

// Package logger provides support for logging to stdout and stderr.
// Log entries will be logged with format: $timestamp $hostname [$pid]: $severity $message.
package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// LogFormatter is used to format log entry.
type LogFormatter struct{}

// Format formats a given log entry, returns byte slice and error.
func (c *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Format(time.RFC3339)

	level := strings.ToUpper(entry.Level.String())
	if len(level) == 4 {
		level = " " + level
	}

	return []byte(fmt.Sprintf(
		"[%s #%d] %s -- : %s\n",
		timestamp,
		os.Getpid(),
		level,
		entry.Message)), nil
}

func init() {
	log.SetFormatter(&LogFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)
}

// GetLevel get the log level string.
func GetLevel() string {
	return log.GetLevel().String()
}

// SetLevel sets the log level. Valid levels are "debug", "info", "warn", "error", and "fatal".
func SetLevel(level string) {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		Fatal(fmt.Sprintf(`not a valid level: "%s"`, level))
	}
	log.SetLevel(lvl)
}

// Debug logs a message with severity DEBUG.
func Debug(format string, v ...interface{}) {
	log.Debug(fmt.Sprintf(format, v...))
}

// Info logs a message with severity INFO.
func Info(format string, v ...interface{}) {
	log.Info(fmt.Sprintf(format, v...))
}

// Warn logs a message with severity WARNING.
func Warn(format string, v ...interface{}) {
	log.Warn(fmt.Sprintf(format, v...))
}

// Error logs a message with severity ERROR.
func Error(format string, v ...interface{}) {
	log.Error(fmt.Sprintf(format, v...))
}

// Fatal logs a message with severity ERROR followed by a call to os.Exit().
func Fatal(format string, v ...interface{}) {
	log.Fatal(fmt.Sprintf(format, v...))
}
