package commands

import (
	"context"

	"github.com/pengsrc/go-shared/pid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qscamel/config"
	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/model"
)

// CleanCmd will clean all finished task.
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean all finished task.",
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
}
