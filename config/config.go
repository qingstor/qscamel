package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/utils"
)

// Config stores all config value.
type Config struct {
	ThreadNum int `yaml:"thread_num"`

	LogFile      string `yaml:"log_file"`
	LogLevel     string `yaml:"log_level"`
	PIDFile      string `yaml:"pid_file"`
	DatabaseFile string `yaml:"database_file"`
}

// New will create a new Config.
func New() (*Config, error) {
	return &Config{}, nil
}

// LoadFromFilePath will load config from specific file path.
func (c *Config) LoadFromFilePath(filePath string) (err error) {
	filePath, err = utils.Expand(filePath)
	if err != nil {
		return
	}

	f, err := os.OpenFile(filePath, os.O_RDWR, 0600)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
		f, err = utils.CreateFile(filePath)
		if err != nil {
			return
		}
		_, err = f.WriteString(constants.DefaultConfigContent)
		if err != nil {
			return
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			return
		}
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return c.LoadFromContent(content)
}

// LoadFromContent will load config from file content.
func (c *Config) LoadFromContent(content []byte) error {
	return yaml.Unmarshal(content, c)
}

// Check will check whether the config is vaild.
func (c *Config) Check() (err error) {
	// Check thread number.
	if c.ThreadNum == 0 {
		c.ThreadNum = constants.DefaultThreadNum
	}

	// Check log file.
	if c.LogFile == "" {
		c.LogFile = constants.LogPath
	}
	c.LogFile, err = utils.Expand(c.LogFile)
	if err != nil {
		return
	}
	_, err = os.OpenFile(c.LogFile, os.O_RDWR, 0600)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
		_, err = utils.CreateFile(c.LogFile)
		if err != nil {
			return
		}
	}

	// Check database file.
	if c.DatabaseFile == "" {
		c.DatabaseFile = constants.DatabasePath
	}
	c.DatabaseFile, err = utils.Expand(c.DatabaseFile)
	if err != nil {
		return
	}
	_, err = os.OpenFile(c.DatabaseFile, os.O_RDWR, 0600)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
		_, err = utils.CreateFile(c.DatabaseFile)
		if err != nil {
			return
		}
	}

	return nil
}
