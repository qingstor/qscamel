package contexts

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/yunify/qscamel/config"
	"github.com/yunify/qscamel/db"
	"github.com/yunify/qscamel/utils"
)

var (
	// DB holds the database connection.
	DB *db.Database
	// Config stores the current config.
	Config *config.Config
	// Client stores the http client used in qscamel.
	Client *http.Client
	// Proxy
	Proxy *url.URL
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
	f := &lumberjack.Logger{
		Filename:  c.LogFile,
		MaxSize:   1024,
		LocalTime: true,
		Compress:  true,
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, f))

	// Setup Bolt.
	DB, err = db.NewDB(&db.DatabaseOptions{
		Address: c.DatabaseFile,
	})
	if err != nil {
		return
	}

	var proxy *url.URL
	if c.Proxy != "" {
		proxy, err = url.Parse(c.Proxy)
		if err != nil {
			return
		}
	}
	Proxy = proxy

	// Setup http client.
	Client = &http.Client{
		// We do not use the timeout in http client,
		// because this timeout is for the whole http body read/write,
		// it's unsuitable for various length of files and network condition.
		// We provide a wrapper in utils/conn.go of net.Dialer to make io timeout
		// to the http connection for individual buffer I/O operation,
		Timeout:   0,
		Transport: NewTransportWithDialContext(c, proxy, utils.DefaultDialer),
	}

	return nil
}

func NewTransportWithDialContext(c *config.Config, proxy *url.URL, dialer *utils.Dialer) *http.Transport {
	return &http.Transport{
		DialContext: dialer.DialContext,
		// Client will be used in both source and destination.
		MaxIdleConns: c.Concurrency * 2,
		// Max idle conns should be config's concurrency.
		MaxIdleConnsPerHost: c.Concurrency,

		IdleConnTimeout:       time.Second * 20,
		TLSHandshakeTimeout:   time.Second * 10,
		ExpectContinueTimeout: time.Second * 2,

		Proxy: http.ProxyURL(proxy),
	}
}
