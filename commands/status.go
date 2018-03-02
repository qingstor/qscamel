package commands

import (
	"context"

	"github.com/pengsrc/go-shared/pid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qscamel/config"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/model"
)

// StatusCmd will show current task status.
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current task status.",
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
}
