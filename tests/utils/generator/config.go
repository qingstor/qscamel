package generator

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/config"
)

// ConfigContentfmt is alias of Config struct
type ConfigContentfmt config.Config

// TaskContentfmt is struct that will
// be serialized to yaml format
type TaskContentfmt struct {
	Tasktype string `yaml:"type"`
	Src      struct {
		FSType  string      `yaml:"type"`
		Path    string      `yaml:"path"`
		Options interface{} `yaml:"options,omitempty"`
	} `yaml:"source"`
	Dst struct {
		FSType  string      `yaml:"type"`
		Path    string      `yaml:"path"`
		Options interface{} `yaml:"options,omitempty"`
	} `yaml:"destination"`
}

func confAssign(dir string) *ConfigContentfmt {
	return &ConfigContentfmt{
		0, dir + "/qscamel.log",
		"info", dir + "/qscamel.pid",
		dir + "/db",
	}
}

func taskAssign(dir, tskType, srcFs, dstFs string,
	srcOpt interface{}, dstOpt interface{}) *TaskContentfmt {
	t := TaskContentfmt{}
	t.Tasktype = tskType
	t.Src.Path = dir + "/src"
	t.Src.FSType = srcFs
	t.Src.Options = srcOpt
	t.Dst.Path = dir + "/dst"
	t.Dst.FSType = dstFs
	t.Dst.Options = dstOpt
	return &t
}

// CreateTestConfigYaml create config yaml file for test
// in the `dir` directory, and return the config file
// name if there are no errors.
func CreateTestConfigYaml(dir string) (string, error) {
	confFile, err := ioutil.TempFile(dir, "config*.yaml")
	if err != nil {
		return "", err
	}
	defer confFile.Close()
	confContent, err := yaml.Marshal(confAssign(dir))
	if err != nil {
		return "", err
	}
	if _, err := confFile.Write(confContent); err != nil {
		return "", err
	}
	return confFile.Name(), nil
}

// CreateTestTaskYaml creat task yaml file for test
// in the `dir` directory, and return the task file name
// if there are no errors
func CreateTestTaskYaml(dir, tskType, srcFs, dstFs string,
	srcOpt, dstOpt interface{}) (string, error) {
	taskFile, err := ioutil.TempFile(dir, "task*.yaml")
	defer taskFile.Close()
	if err != nil {
		return "", err
	}
	taskContent, err := yaml.Marshal(taskAssign(dir, tskType, srcFs, dstFs, srcOpt, dstOpt))
	if err != nil {
		return "", err
	}
	if _, err := taskFile.Write(taskContent); err != nil {
		return "", err
	}
	return taskFile.Name(), nil
}
