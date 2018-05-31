package commands

import (
	"context"
	"os"
	"os/signal"

	"github.com/pengsrc/go-shared/pid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qscamel/config"
	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/migrate"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// RunCmd will provide run command for qscamel.
var RunCmd = &cobra.Command{
	Use:   "run [task name or task path]",
	Short: "Create or resume a task",
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

		ctx := context.Background()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, os.Kill)
		go func() {
			sig := <-sigs
			logrus.Infof("%v signal received, exit for now.", sig)

			contexts.DB.Close()

			os.Exit(0)
		}()

		// Load and check task.
		t, err := model.LoadTask(args[0])
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

		err = migrate.Execute(ctx)
		if err != nil {
			logrus.Errorf("Migrate failed for %v.", err)
		}
	},
}
