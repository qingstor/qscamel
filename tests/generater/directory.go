package generater

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/yunify/qscamel/utils"
)

// CleanTestTempFile will clean the temp file which created
// by corresponded task.
func CleanTestTempFile(fmap *map[string]string) error {
	if err := os.RemoveAll((*fmap)["dir"]); err != nil {
		return err
	}
	return nil
}

// CreateTestConfigFile create the config file for one test
// it return a mapping of some configuration of the test
// "dir" is the base directory path of the test
// "config" is the config file path (point to database path
// , pid file path etc.)
// "task" is the task file path for run a task on random path
func CreateTestConfigFile(tskType, srcFs, dstFs string,
	srcOpt, dstOpt interface{}, p bool) (*map[string]string, error) {
	fileMap := make(map[string]string)

	// create temp directory
	dir, err := ioutil.TempDir("", "qscamel")
	if err != nil {
		return nil, err
	}
	fileMap["dir"] = dir

	// create a temp config file
	confName, err := CreateTestConfigYaml(dir)
	if err != nil {
		return nil, err
	}
	fileMap["config"] = confName

	// create a temp task file
	taskName, err := CreateTestTaskYaml(dir, tskType, srcFs, dstFs, srcOpt, dstOpt)
	if err != nil {
		return nil, err
	}
	fileMap["task"] = taskName
	fileMap["name"] = extractTaskName(taskName)

	if p {
		fmt.Println("create temp dir at", dir)
		fmt.Println("create temp config file at ", confName)
		fmt.Println("create temp task file at ", taskName)
	}
	return &fileMap, nil
}

// CreateTestDefaultFile will be used to generate
// task file, but the config file will be yield
// by qscamel itself
func CreateTestDefaultFile(tskType, srcFs, dstFs string,
	srcOpt, dstOpt interface{}, p bool) (*map[string]string, error) {
	fileMap := make(map[string]string)
	home, err := utils.Dir()
	if err != nil {
		return nil, err
	}
	fileMap["dir"] = home + "/.qscamel"
	if err := os.MkdirAll(fileMap["dir"], 0700); err != nil {
		fmt.Println(err)
		return nil, err
	}
	taskname, err := CreateTestTaskYaml(fileMap["dir"], tskType, srcFs, dstFs, srcOpt, dstOpt)
	if err != nil {
		return nil, err
	}
	fileMap["task"] = taskname
	fileMap["name"] = extractTaskName(taskname)
	return &fileMap, nil
}

func extractTaskName(pn string) string {
	// extract the taskname(taskXXXXX) from task file path(/tmp/qscamelXXXXX/taskXXXX.yaml) .
	_, taskName := path.Split(pn)
	runName := strings.Split(taskName, ".")
	return runName[0]
}

// CreateLocalSrcTestRandDirFile generate the random name directory and file in
// the base directory in the `fmap`. it create `filePerDir` numbers file and
// `dirPerDir` numbers directory in every directory, and the file size is `fileSize`
// `dirDepth` point to the directory depth to generate(advised depth is `2`).
func CreateLocalSrcTestRandDirFile(fmap *map[string]string, filePerDir int, dirPerDir int,
	fileSize int64, dirDepth int) error {
	err := os.MkdirAll((*fmap)["dir"]+"/src", 0755)
	if err != nil {
		return err
	}
	(*fmap)["src"] = (*fmap)["dir"] + "/src"

	chsz := seriesSum(dirPerDir, dirDepth)
	subchsz := seriesSum(dirPerDir, dirDepth-1)
	dirch := make(chan string, chsz)
	done := make(chan error, 0)

	// generate create directory recursively task for goroutine
	if chsz >= 1 {
		dirch <- (*fmap)["src"]
	}

	go func() {
		for i := 0; i < chsz && subchsz > 0; i++ {
			select {
			case onePath := <-dirch:
				if err := CreateTestRandomFile(filePerDir, fileSize, onePath); err != nil {
					done <- err
				}
				if i >= subchsz {
					continue
				}
				if err := CreateTestSubDirectory(dirch, dirPerDir, onePath); err != nil {
					done <- err
				}
			default:
				done <- nil
			}
		}
		done <- nil
	}()
	if dirDepth == 1 {
		if err := CreateTestRandomFile(filePerDir, fileSize, (*fmap)["src"]); err != nil {
			return err
		}
	}

	return <-done
}

// CreateTestSubDirectory create `dirPerDir` number of directory in
// `dir` directory
func CreateTestSubDirectory(dirch chan string, dirPerDir int, dir string) error {
	for ; dirPerDir > 0; dirPerDir-- {
		name, err := ioutil.TempDir(dir, "DIR")
		if err != nil {
			return err
		}
		dirch <- name
	}
	return nil
}

// CreateLocalDstDir create the destination directory
// in the local machine
func CreateLocalDstDir(fmap *map[string]string) error {
	err := os.MkdirAll((*fmap)["dir"]+"/dst", 0755)
	if err != nil {
		return err
	}
	(*fmap)["dst"] = (*fmap)["dir"] + "/dst"
	return nil
}

// CreateLocalSrcDir create the source directory
// in the local machine
func CreateLocalSrcDir(fmap *map[string]string) error {
	err := os.MkdirAll((*fmap)["dir"]+"/src", 0755)
	if err != nil {
		return err
	}
	(*fmap)["src"] = (*fmap)["dir"] + "/src"
	return nil
}
