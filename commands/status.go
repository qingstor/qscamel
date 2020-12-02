package commands

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qscamel/model"
)

// StatusCmd will show current task status.
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current task status.",
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return initContext(cmd.Flag("config").Value.String())
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Start show status.
		logrus.Infof("Show status started.")

		ctx := context.Background()
		t, err := model.ListTask(ctx)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.Printf("There are %d tasks totally.", len(t))
		for _, v := range t {
			logrus.Printf("Task: %s, Status: %s", v.Name, v.Status)
		}
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return cleanUp()
	},
}
