package checker

import (
	"testing"

	"github.com/yunify/qscamel/tests/executer"
	"github.com/yunify/qscamel/tests/generater"
)

func TestTaskRunCopy(t *testing.T) {
	// env set
	fileMap, err := generater.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, false)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 4, 1, generater.MB*2, 2, true)
	err = generater.CreateLocalDstDir(fileMap)
	if err != nil {
		t.Fatal(err)
	}

	// run command
	if err = executer.Execute(fileMap, "run"); err != nil {
		t.Fatal(err)
	}
	// check running ouput
	if err := executer.CheckOutput(fileMap,
		"Start copying single object [A-Z0-9]*/TESTFILE\\d+.camel", 8, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task [a-z0-9]* has been finished", 1, true); err != nil {
		t.Fatal(err)
	}
}

func TestTaskDelete(t *testing.T) {
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
	// run command
	if err = executer.Execute(fileMap, "run"); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task [a-z0-9]* has been finished", 1, true); err != nil {
		t.Fatal(err)
	}
	// delete command
	(*fileMap)["delname"] = (*fileMap)["name"]
	if err := executer.Execute(fileMap, "delete"); err != nil {
		t.Fatal(err)
	}
	// check delete output
	if err := executer.CheckOutput(fileMap,
		"Task [a-z0-9]* has been deleted", 1, true); err != nil {
		t.Fatal(err)
	}

}

func TestTaskStatus(t *testing.T) {
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
	// run command
	if err = executer.Execute(fileMap, "run"); err != nil {
		t.Fatal(err)
	}
	// status command
	if err = executer.Execute(fileMap, "status"); err != nil {
		t.Fatal(err)
	}
	// check status output
	if err := executer.CheckOutput(fileMap,
		"Show status started", 1, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"There are 1 tasks totally", 1, true); err != nil {
		t.Fatal(err)
	}

}

func TestTaskClean(t *testing.T) {
	fileMap, err := generater.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, false)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 1, 0, generater.MB*2, 1, true)
	err = generater.CreateLocalDstDir(fileMap)
	if err != nil {
		t.Fatal(err)
	}
	// run command
	if err = executer.Execute(fileMap, "run"); err != nil {
		t.Fatal(err)
	}
	// clean command
	if err = executer.Execute(fileMap, "clean"); err != nil {
		t.Fatal(err)
	}
	// check clean output
	if err := executer.CheckOutput(fileMap,
		"Clean started", 1, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task [a-z0-9]* has been cleaned", 1, true); err != nil {
		t.Fatal(err)
	}
}
