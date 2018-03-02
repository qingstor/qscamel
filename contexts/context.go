package contexts

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"

	"github.com/yunify/qscamel/config"
	"github.com/yunify/qscamel/db"
)

var (
	// DB holds the database connection.
	DB *db.Database
	// Config stores the current config.
	Config *config.Config
)

// SetupContexts will set contexts.
func SetupContexts(c *config.Config) (err error) {
	// Setup config.
	Config = c

	// Setup logger.
	// Set log level.
	lvl, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		lvl = logrus.ErrorLevel
	}
	logrus.SetLevel(lvl)
	// Set formatter.
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	// Set output.
	f, err := os.OpenFile(c.LogFile, os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, f))

	// Setup Bolt.
	DB, err = db.NewDB(&db.DatabaseOptions{
		Address: c.DatabaseFile,
	})
	if err != nil {
		return
	}
	err = DB.Init()
	if err != nil {
		return
	}

	return nil
}
