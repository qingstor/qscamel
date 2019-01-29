package checker

import (
	"github.com/yunify/qscamel/tests/executer"
	"testing"

	"github.com/yunify/qscamel/tests/generater"
)

func TestEmptyDirectory(t *testing.T) {
	// env set
	fileMap, err := generater.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, true)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 0, 0, generater.MB*2, 1)
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
		"Task task[0-9]* migrate started", 1, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* has been finished", 1, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutputUnexpect(fileMap, "Job /DIR[0-9]* listed", true); err != nil {
		t.Fatal(err)
	}

}

func TestOneDirectory(t *testing.T) {
	fileMap, err := generater.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, true)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 0, 1, generater.MB*2, 2)
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
		"Task task[0-9]* migrate started", 1, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* has been finished", 1, true); err != nil {
		t.Fatal(err)
	}

}

func TestDeepDirectory(t *testing.T) {
	fileMap, err := generater.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, true)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 0, 1, generater.MB*2, 257)
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
		"Task task[0-9]* migrate started", 1, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* has been finished", 1, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Job (/DIR[0-9]*)* listed", 256, false); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutputUnexpect(fileMap,
		"Start copying single object [A-Z0-9]*/TESTFILE\\d+.camel", true); err != nil {
		t.Fatal(err)
	}
}

func TestManyDirectory(t *testing.T) {
	fileMap, err := generater.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, true)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 0, 3, 0, 3)
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
		"Task task[0-9]* migrate started", 1, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* has been finished", 1, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Job (/DIR[0-9]*)* listed", 36, false); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutputUnexpect(fileMap,
		"Start copying single object [A-Z0-9]*/TESTFILE\\d+.camel", true); err != nil {
		t.Fatal(err)
	}
}
