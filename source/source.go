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

package source

import (
	"fmt"
	"time"

	log "github.com/frostyplanet/logrus"

	"github.com/yunify/qscamel/record"
	"github.com/yunify/qscamel/source/aliyun"
	"github.com/yunify/qscamel/source/file"
	"github.com/yunify/qscamel/source/qingstor"
	"github.com/yunify/qscamel/source/qiniu"
	"github.com/yunify/qscamel/source/s3"
	"github.com/yunify/qscamel/source/upyun"
)

const (
	sourceTypeFile     = "file"
	sourceTypeAliyun   = "aliyun"
	sourceTypeQiniu    = "qiniu"
	sourceTypeS3       = "s3"
	sourceTypeUpyun    = "upyun"
	sourceTypeQingstor = "qingstor"
)

// MigrateSource is interface of migrating source specified by '--src-type' flag.
type MigrateSource interface {
	// GetSourceSites reads source sites list from source.
	GetSourceSites(threadNum int, logger *log.Logger, recorder *record.Recorder,
	) (sourceSites []string, objectNames []string, skipped []string, done bool, err error)

	// GetSourceSiteInfo returns last modified time of the given source site.
	GetSourceSiteInfo(sourceSite string) (lastModified time.Time, err error)
}

// InstantiateMigrateSource creates concrete source instance.
func InstantiateMigrateSource(sourceType, specificSource, zone, accessKey, secreteKey string,
) (MigrateSource, error) {
	switch sourceType {
	case sourceTypeFile:
		return file.NewSourceFile(specificSource)
	case sourceTypeAliyun:
		return aliyun.NewSourceAliyun(specificSource, zone, accessKey, secreteKey)
	case sourceTypeQiniu:
		return qiniu.NewSourceQiniu(specificSource, accessKey, secreteKey)
	case sourceTypeS3:
		return s3.NewSourceS3(specificSource, zone, accessKey, secreteKey)
	case sourceTypeUpyun:
		return upyun.NewSourceUpyun(specificSource, zone, accessKey, secreteKey)
	case sourceTypeQingstor:
		return qingstor.NewSourceQingstor(specificSource, zone, accessKey, secreteKey)
	default:
		return nil, fmt.Errorf("Unsupported source type: %s", sourceType)
	}
}
