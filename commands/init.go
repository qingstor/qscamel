package commands

import (
	"github.com/pengsrc/go-shared/pid"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/config"
	"github.com/yunify/qscamel/contexts"
)

func init() {
	RunCmd.Flags().StringVarP(&taskPath, "task", "t", "", "task path")
}

func initContext(configFile string) error {
	c := &config.Config{}
	if err := c.LoadFromFilePath(configFile); err != nil {
		logrus.Errorf("Load config from %s failed for %v.", configFile, err)
		return err
	}

	// Check config.
	if err := c.Check(); err != nil {
		logrus.Errorf("Config check failed for %v.", err)
		return err
	}

	// Create PID file.
	if pidfile := c.PIDFile; pidfile != "" {
		p, err := pid.New(pidfile)
		if err != nil {
			logrus.Errorf("PID create failed for %v.", err)
			return err
		}
		defer func() {
			err = p.Remove()
			if err != nil {
				logrus.Errorf("PID remove failed for %v.", err)
			}
		}()
	}

	// Setup contexts.
	if err := contexts.SetupContexts(c); err != nil {
		logrus.Errorf("Contexts setup failed for %v.", err)
		return err
	}
	return nil
}

func cleanUp() error {
	if contexts.DB != nil {
		contexts.DB.Close()
	}
	return nil
}
