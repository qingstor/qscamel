package generater

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"testing"
)


var (
	ConfigContentfmtForTest =
`
concurrency: 0
log_file: %s/qscamel.log
log_level: info
pid_file: %s/qscamel.pid
database_file: %s/db`
	TaskContentfmtForTest =
`
type: copy
source:
  type: fs
  path: %s/src
destination:
  type: fs
  path: %s/dst
`)




const B int64 = 1
const KB int64 = 1 * 1024
const MB int64 = KB * 1024
const GB int64  = MB * 1024

// CleanTestTempFile will clean the temp file which created
// by corresponded task.
func CleanTestTempFile(t *testing.T, fmap map[string]string) {
	if err := os.RemoveAll(fmap["dir"]); err != nil {
		t.Fatal(err)
	}
}

// CreateTestConfigFile create the config file for one test
// it return a mapping of some configuration of the test
// "dir" is the base directory path of the test
// "config" is the config file path (point to database path
// , pid file path etc.)
// "task" is the task file path for run a task on random path
func CreateTestConfigFile(t *testing.T) map[string]string{
	fileMap := make(map[string] string)

	// create temp directory
	dir, err := ioutil.TempDir("", "qscamel")
	if err != nil {
		t.Fatal(err)
	}
	fileMap["dir"] = dir
	fmt.Println("create temp dir at", dir)

	// create a temp config file
	confContent := fmt.Sprintf(ConfigContentfmtForTest, dir, dir, dir)
	confFile, err := ioutil.TempFile(dir, "config*.yaml")
	if err != nil {
		Fatal(t, fileMap, err)
	}
	if _, err := confFile.WriteString(confContent); err != nil {
		t.Fatal(err)
	}
	fileMap["config"] = confFile.Name()
	fmt.Println("create temp config file at ", confFile.Name())
	// create a temp task file
	taskContent := fmt.Sprintf(TaskContentfmtForTest, dir, dir)
	taskFile, err := ioutil.TempFile(dir, "task*.yaml")
	if err != nil {
		Fatal(t, fileMap, err)
	}
	if _, err := taskFile.WriteString(taskContent); err != nil {
		Fatal(t, fileMap, err)
	}
	fileMap["task"] = taskFile.Name()
	fmt.Println("create temp task file at ", taskFile.Name())
	return fileMap
}

func Fatal(t *testing.T, fmap map[string]string, err error){
	CleanTestTempFile(t, fmap)
	t.Fatal(err)
}


// CreateTestRandDirFile generate the random name directory and file in
// the base directory in the `fmap`. it first creat two directory, and the
// name is "src" and "dst" respective, it create `filePerDir` file and
// `dirPerDir` directory in every directory, and the file size is `fileSize`
// `dirDepth` point to the directory depth to generate(best `2`).
func CreateTestRandDirFile(t *testing.T, fmap map[string]string,
	filePerDir int, dirPerDir int, fileSize int64, dirDepth int, isRandom bool) {
	err := os.MkdirAll(fmap["dir"] + "/src",0755)
	err = os.MkdirAll(fmap["dir"] + "/dst", 0755)
	if err != nil {
		Fatal(t, fmap, err)
	}
	fmap["src"] = fmap["dir"]+"/src"
	fmap["dst"] = fmap["dir"]+"/dst"

	chsz := seriesSum(dirPerDir, dirDepth)
	subchsz := seriesSum(dirPerDir, dirDepth-1)
	dirch := make(chan string, chsz)
	done := make(chan int, 0)
	if err := CreateTestSubDirectory(dirch, dirPerDir, fmap["src"]); err != nil {
		Fatal(t, fmap, err)
	}
	//fmt.Println(chsz, subchsz)
	go func() {
		for i := 0;i < chsz; i++{
			if path, ok := <-dirch; ok != false{
				fmt.Println("create temp directory", path)

				if err := CreateTestRandomFile(filePerDir, fileSize, path); err != nil {
					Fatal(t, fmap, err)
				}
				if i >= subchsz {
					continue
				}
				if err := CreateTestSubDirectory(dirch, dirPerDir, path); err != nil {
					Fatal(t, fmap, err)
				}
			}
		}
		done <-1
	}()

	if err := CreateTestRandomFile(filePerDir, fileSize, fmap["src"]); err != nil {
		Fatal(t, fmap, err)
	}
	<-done


}

// CreateTestRandomFile create the `filePerDir` number of random file
// in the `dir` directory.
func CreateTestRandomFile(filePerDir int, fileSize int64, dir string) error {
	for i := 0; i < filePerDir; i++ {
		name, err := ioutil.TempFile(dir+"/", "TestFile*.camel")
		if err == nil {
			content, err := CreateRandomByteStream(fileSize)
			if err == nil {
				_, err := name.Write(content)
				if err == nil {
					continue
				}
			}
		}
		return err
	}
	return nil
}

// CreateTestSubDirectory create the `dirPerDir` number of directory in
// `dir` directory
func CreateTestSubDirectory(dirch chan string, dirPerDir int, dir string) error {
	for ; dirPerDir >0; dirPerDir-- {
		name, err := ioutil.TempDir(dir, "DIR")
		fmt.Println(name)
		if err != nil {
			return err
		}
		dirch <- name
	}
	return nil
}

// CreateRandomByteStream return `size` of random bytes
func CreateRandomByteStream(size int64) ([]byte,error) {
	p := make([]byte, size)
	_, err := rand.Read(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// caculate Geometric series sum
func seriesSum(dirnum, depth int) int {
	return dirnum * int(((1 - math.Pow(float64(dirnum), float64(depth))) /float64( 1 - dirnum)))
}