package utils

import (
	"io/ioutil"
	"testing"

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
func CreateTestConfigYaml(t testing.TB, dir string) string {
	confFile, err := ioutil.TempFile(dir, "config*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer confFile.Close()
	confContent, err := yaml.Marshal(confAssign(dir))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := confFile.Write(confContent); err != nil {
		t.Fatal(err)
	}
	return confFile.Name()
}

// CreateTestTaskYaml creat task yaml file for test
// in the `dir` directory, and return the task file name
// if there are no errors
func CreateTestTaskYaml(t testing.TB, dir, tskType, srcFs, dstFs string,
	srcOpt, dstOpt interface{}) string {
	taskFile, err := ioutil.TempFile(dir, "task*.yaml")
	defer taskFile.Close()
	if err != nil {
		t.Fatal(err)
	}
	taskContent, err := yaml.Marshal(taskAssign(dir, tskType, srcFs, dstFs, srcOpt, dstOpt))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := taskFile.Write(taskContent); err != nil {
		t.Fatal(err)
	}
	return taskFile.Name()
}
