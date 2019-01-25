package checker

import (
	"testing"
	"time"

	"github.com/yunify/qscamel/tests/dbtester"
	"github.com/yunify/qscamel/tests/generater"
)

func TestTaskRunCopy(t *testing.T) {
	// env set
	fileMap, err := generater.CreateTestConfigFile("copy", "fs", "fs", nil, nil)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(fileMap, 10, 3, generater.MB*2, 1, true)
	err = generater.CreateLocalDstDir(fileMap)
	time.Sleep(100 * time.Second)
	if err != nil {
		t.Fatal(err)
	}

	// check DB
	if err := dbtester.CheckDBEmpty((*fileMap)["dir"]); err != nil {
		t.Fatal(err)
	}
	// check Ouput
}

func TestTaskDelete(t *testing.T) {

}

func TestTaskStatus(t *testing.T) {

}

func TestTaskClean(t *testing.T) {

}
