package checker

import (
	"testing"

	"github.com/yunify/qscamel/tests/executer"
	"github.com/yunify/qscamel/tests/generater"
)

func TestEmptyFile(t *testing.T) {
	// env set
	fileMap, err := generater.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, true)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 3, 0, 0, 1)
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

	// check file equal
	if err := executer.CheckDirectroyEqual(fileMap); err != nil {
		t.Fatal(err)
	}

}

// Bugfix : BigFile connot be gennerated
func testBigFile(t *testing.T) {
	// env set
	fileMap, err := generater.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, true)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 1, 0, 3*generater.GB, 1)
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

	// check file equal
	if err := executer.CheckDirectroyEqual(fileMap); err != nil {
		t.Fatal(err)
	}
}

func TestManyFile(t *testing.T) {
	// env set
	fileMap, err := generater.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, true)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 100000, 0, 1, 1)
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
		"Start copying single object [A-Z0-9/]*/TESTFILE\\d+.camel", 100000, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* has been finished", 1, true); err != nil {
		t.Fatal(err)
	}

	// check file equal
	if err := executer.CheckDirectroyEqual(fileMap); err != nil {
		t.Fatal(err)
	}
}

func TestDeepFile(t *testing.T) {
	// env set
	fileMap, err := generater.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, true)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 1, 1, 3*generater.MB, 10)
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
		"Start copying single object [A-Z0-9/]*/TESTFILE\\d+.camel", 10, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* has been finished", 1, true); err != nil {
		t.Fatal(err)
	}

	// check file equal
	if err := executer.CheckDirectroyEqual(fileMap); err != nil {
		t.Fatal(err)
	}
}

func TestMutiDirAndFile(t *testing.T) {
	// env set
	fileMap, err := generater.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, true)
	defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 10, 3, 2*generater.MB, 3)
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
		"Start copying single object [A-Z0-9/]*/TESTFILE\\d+.camel", 416, true); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* has been finished", 1, true); err != nil {
		t.Fatal(err)
	}

	// check file equal
	if err := executer.CheckDirectroyEqual(fileMap); err != nil {
		t.Fatal(err)
	}
}
