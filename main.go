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

package main

import (
	"github.com/spf13/cobra"

	"github.com/yunify/qscamel/commands"
	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/utils"
)

var application = &cobra.Command{
	Use:   constants.Name,
	Short: constants.ShortDescription,
	Long:  constants.LongDescription,
}

var (
	configPath string
)

func init() {
	// Add version command.
	application.AddCommand(commands.VersionCmd)
	// Add run command.
	application.AddCommand(commands.RunCmd)
	// Add delete command.
	application.AddCommand(commands.DeleteCmd)
	// Add clean command.
	application.AddCommand(commands.CleanCmd)
	// Add status command.
	application.AddCommand(commands.StatusCmd)

	// Add config flag which can be used in all sub commands.
	application.PersistentFlags().StringVarP(&configPath, "config", "c", constants.ConfigPath, "config path")
	application.MarkFlagRequired("config")
}

func main() {
	utils.CheckError(application.Execute)
}
