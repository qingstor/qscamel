package commands

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// DeleteCmd will provide delete command for qscamel.
var DeleteCmd = &cobra.Command{
	Use:   "delete [task name]",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return initContext(cmd.Flag("config").Value.String())
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Start delete.
		ctx := context.Background()
		// Load and check task.
		t, err := model.GetTaskByName(ctx, args[0])
		if err != nil {
			logrus.Panicf("Task load failed for %v.", err)
			return
		}
		if t == nil {
			logrus.Errorf("Task %s is not exist.", args[0])
			return
		}

		ctx = utils.NewTaskContext(ctx, t.Name)

		// Start delete.
		logrus.Infof("Task %s delete started.", t.Name)

		err = model.DeleteTask(ctx)
		if err != nil {
			logrus.Panicf("Delete task for %v.", err)
		}

		logrus.Infof("Task %s has been deleted.", t.Name)
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return cleanUp()
	},
}
