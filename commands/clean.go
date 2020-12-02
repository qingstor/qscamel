package commands

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// CleanCmd will clean all finished task.
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean all finished task.",
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return initContext(cmd.Flag("config").Value.String())
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Start clean.
		logrus.Infof("Clean started.")

		ctx := context.Background()
		t, err := model.ListTask(ctx)
		if err != nil {
			logrus.Panic(err)
		}
		for _, v := range t {
			if v.Status != constants.TaskStatusFinished {
				continue
			}

			err = model.DeleteTaskByName(ctx, v.Name)
			if err != nil {
				logrus.Panic(err)
			}

			logrus.Infof("Task %s has been cleaned.", v.Name)
		}
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return cleanUp()
	},
}
