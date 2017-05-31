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
	"errors"
	"fmt"
	"net/http"
	"os"

	log "github.com/frostyplanet/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qingstor-sdk-go/config"
	"github.com/yunify/qscamel/metadata"
	"github.com/yunify/qscamel/record"
	"github.com/yunify/qscamel/source"
	"github.com/yunify/qscamel/utils"
)

const (
	// DefaultThreadNum is default num of objects being migrated concurrently.
	DefaultThreadNum = 10

	// MaxThreadNum is max num of objects being migrated concurrently.
	MaxThreadNum = 100
)

// Context holds runtime context for Migrate function.
type Context struct {
	Source     source.MigrateSource
	SourceType string

	QSConfig     *config.Config
	QSBucketName string

	PrintVersion     bool
	Overwrite        bool
	IgnoreExisting   bool
	IgnoreUnmodified bool
	DryRun           bool
	ThreadNum        int
	Logger           *log.Logger

	Recorder *record.Recorder
}

var (
	ctx = &Context{
		IgnoreUnmodified: true,
		Logger:           utils.GetLogger(),
		QSConfig:         &config.Config{Connection: &http.Client{}},
	}
	specificSource        string
	sourceZone            string
	sourceAccessKeyID     string
	sourceSecretAccessKey string
	configPath            string
	logPath               string
	description           string
)

// Execute parses flags, checks required input and calls Migrate function.
func Execute() {
	cmd := &cobra.Command{
		PreRunE: checkFlags,
		Run: func(cmd *cobra.Command, args []string) {
			if ctx.PrintVersion {
				fmt.Printf(
					"qscamel version %s\n", metadata.Version,
				)
				return
			}
			completed, failed, skipped, err := Migrate(ctx)
			utils.CheckErrorForExit(err)
			utils.LogResult(completed, "completed", ctx.Logger, ctx.DryRun)
			utils.LogResult(failed, "failed", ctx.Logger, ctx.DryRun)
			utils.LogResult(skipped, "skipped", ctx.Logger, ctx.DryRun)
		},
	}
	defineFlags(cmd)
	err := cmd.Execute()
	utils.CheckErrorForExit(err)
}

func defineFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(
		&ctx.SourceType, "src-type", "t", "",
		`Specify source type, support "file" and other object storage platform`,
	)
	cmd.Flags().StringVarP(
		&specificSource, "src", "s", "",
		`Specify migration source. If --src-type is "file", --src specifies source list file. Otherwise, --src specifies source bucket name.`,
	)
	cmd.Flags().StringVarP(
		&sourceZone, "src-zone", "z", "",
		"Specify source zone for object storage type source",
	)
	cmd.Flags().StringVarP(
		&sourceAccessKeyID, "src-access-key", "a", "",
		`Specify source access_key_id for object storage type source. If --src-type is "upyun", "src-access-key" specifies the name of upyun operator.`,
	)
	cmd.Flags().StringVarP(
		&sourceSecretAccessKey, "src-secret-key", "S", "",
		`Specify source secret_access_key for object storage type source. If --src-type is "upyun", "src-secret-key" specifies the password of upyun operator.`,
	)
	cmd.Flags().StringVarP(
		&ctx.QSBucketName, "bucket", "b", "",
		"Specify QingStor bucket",
	)
	cmd.Flags().StringVarP(
		&configPath, "config", "c", config.DefaultConfigFile,
		"Specify QingStor yaml configuration file",
	)
	cmd.Flags().BoolVarP(
		&ctx.IgnoreExisting, "ignore-existing", "i", false,
		"Ignore existing object",
	)
	cmd.Flags().BoolVarP(
		&ctx.Overwrite, "overwrite", "o", false,
		"Overwrite existing object",
	)
	cmd.Flags().BoolVarP(
		&ctx.DryRun, "dry-run", "n", false,
		"Perform a trial run with no actual migration",
	)
	cmd.Flags().IntVarP(
		&ctx.ThreadNum, "threads", "T", DefaultThreadNum,
		fmt.Sprintf(
			"Specify the number of objects being migrated concurrently (default %d, max %d)", DefaultThreadNum, MaxThreadNum,
		),
	)
	cmd.Flags().StringVarP(
		&logPath, "log-file", "l", "",
		"Specify the path of log file",
	)
	cmd.Flags().StringVarP(
		&description, "description", "d", "",
		"Describe current migration task. This description will be used as record filename for task resuming.",
	)
	cmd.Flags().BoolVarP(
		&ctx.PrintVersion, "version", "v", false,
		"Print the version number of qscamel and exit",
	)
}

func checkFlags(cmd *cobra.Command, args []string) error {
	var err error
	if ctx.PrintVersion {
		return nil
	}

	if specificSource == "" {
		return errors.New("use -s (--src) flag to specify migration source")
	}
	ctx.Source, err = source.InstantiateMigrateSource(
		ctx.SourceType, specificSource, sourceZone, sourceAccessKeyID, sourceSecretAccessKey,
	)
	if err != nil {
		return err
	}

	// Read QingStor config from config file and check authorization variable.
	if configPath == config.DefaultConfigFile {
		configPath = config.GetUserConfigFilePath()
		ctx.Logger.Printf("Use default configuration file %s.", configPath)
	}
	err = ctx.QSConfig.LoadConfigFromFilePath(configPath)
	if err != nil {
		return fmt.Errorf("can't open or parse configuration file %s (%s)", configPath, err.Error())
	}
	if ctx.QSConfig.SecretAccessKey == "" || ctx.QSConfig.AccessKeyID == "" {
		return fmt.Errorf("miss access_key_id or secret_access_key in configuration file %s", configPath)
	}
	if ctx.QSBucketName == "" {
		return errors.New("use -b (--bucket) flag to specify QingStor bucket")
	}

	if ctx.ThreadNum > MaxThreadNum {
		ctx.Logger.Printf("Threads %d is over limit. Use max threads %d.", ctx.ThreadNum, MaxThreadNum)
		ctx.ThreadNum = MaxThreadNum
	}

	// Redirect log output if log file is specified.
	if logPath != "" {
		logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			ctx.Logger.Infof("Write log to %s", logPath)
			ctx.Logger.Out = logFile
		}

	}
	ctx.Logger.Infof("Migration task description: %s", description)

	if description == "" {
		return errors.New("use -d (--description) flag to describe migration task")
	}
	file, err := record.GetRecordFile(description)
	if err != nil {
		return fmt.Errorf("can't create migration record file (%s)", err.Error())
	}
	ctx.Recorder = record.NewRecorder(file)

	return nil
}
