package checker

import (
	"testing"

	"github.com/yunify/qscamel/tests/dbtester"
	"github.com/yunify/qscamel/tests/executer"
	"github.com/yunify/qscamel/tests/generater"
)

func TestTaskRunCopy(t *testing.T) {
	// env set
	fileMap, err := generater.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, false)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 1, 0, generater.MB*2, 1, true)
	err = generater.CreateLocalDstDir(fileMap)
	if err != nil {
		t.Fatal(err)
	}

	// check DB
	if err := dbtester.CheckDBEmpty((*fileMap)["dir"]); err != nil {
		t.Fatal(err)
	}
	// check Ouput
	if err = executer.Execute(fileMap, "run"); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Start copying single object [A-Z0-9]*/TESTFILE\\d+.camel", 1, true); err != nil {
		t.Fatal(err)
	}

	if err := executer.CheckOutput(fileMap,
		"Task [a-z0-9]* has been finished", 1, true); err != nil {
		t.Fatal(err)
	}
}

func TestTaskDelete(t *testing.T) {

}

func TestTaskStatus(t *testing.T) {

}

func TestTaskClean(t *testing.T) {

}
