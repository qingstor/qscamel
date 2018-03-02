package commands

import (
	"context"

	"github.com/pengsrc/go-shared/pid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qscamel/config"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// DeleteCmd will provide delete command for qscamel.
var DeleteCmd = &cobra.Command{
	Use:   "delete [task name]",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := &config.Config{}
		c.LoadFromFilePath(cmd.Flag("config").Value.String())

		// Check config.
		err := c.Check()
		if err != nil {
			logrus.Errorf("Config check failed for %v.", err)
			return
		}

		// Create PID file.
		if pidfile := c.PIDFile; pidfile != "" {
			p, err := pid.New(pidfile)
			if err != nil {
				logrus.Errorf("PID create failed for %v.", err)
				return
			}
			defer func() {
				err = p.Remove()
				if err != nil {
					logrus.Errorf("PID remove failed for %v.", err)
				}
			}()
		}

		// Setup contexts.
		err = contexts.SetupContexts(c)
		if err != nil {
			logrus.Errorf("Contexts setup failed for %v.", err)
			return
		}
		defer contexts.DB.Close()

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
}
