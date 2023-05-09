package commands

import (
	"context"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/migrate"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

var (
	taskPath string
)

// RunCmd will provide run command for qscamel.
var RunCmd = &cobra.Command{
	Use:   "run [task name or task path]",
	Short: "Create or resume a task",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return initContext(cmd.Flag("config").Value.String())
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var closePrint = make(chan struct{}, 1)
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, os.Kill)
		go func() {
			sig := <-sigs
			logrus.Infof("Signal %v received, exit for now.", sig)

			closePrint <- struct{}{}
			migrate.SaveTask()

			cleanUp()
			os.Exit(0)
		}()

		// Load and check task.
		t, err := model.LoadTask(args[0], taskPath)
		if err != nil {
			logrus.Errorf("Task load failed for %v.", err)
			return
		}
		err = t.Check()
		if err != nil {
			logrus.Errorf("Task check failed for %v.", err)
			return
		}

		ctx = utils.NewTaskContext(ctx, t.Name)

		// Start migrate.
		logrus.Infof("Current version: %s.", constants.Version)
		logrus.Infof("Task %s migrate started.", t.Name)

		err = migrate.Execute(ctx, closePrint)
		if err != nil {
			logrus.Errorf("Migrate failed for %v.", err)
		}
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return cleanUp()
	},
}
