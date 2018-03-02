package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yunify/qscamel/constants"
)

// VersionCmd will provide version command for qscamel.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of qscamel",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s\n", constants.Name, constants.Version)
	},
}
