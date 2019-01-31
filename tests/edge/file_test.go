package edge

import (
	"testing"

	"github.com/yunify/qscamel/tests/utils/executer"
	"github.com/yunify/qscamel/tests/utils/generator"
)

var (
	printReg  bool
	printFile bool
)

func init() {
	printReg = false
	printFile = false
}

func TestEmptyFile(t *testing.T) {
	// env set
	fileMap, clean, err := generator.PrepareDefinedTest(3, 0, 0, 1, printFile)
	defer clean(fileMap)
	if err != nil {
		t.Fatal(err)
	}

	// run command
	if err = executer.Execute(fileMap, "run"); err != nil {
		t.Fatal(err)
	}
	// check running ouput
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* migrate started", 1, printReg); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* has been finished", 1, printReg); err != nil {
		t.Fatal(err)
	}

	// check file equal
	if err := executer.CheckDirectroyEqual(fileMap); err != nil {
		t.Fatal(err)
	}
}

// TODO: File Size should be adjusted appropriately
func TestBigFile(t *testing.T) {
	// env set
	fileMap, clean, err := generator.PrepareDefinedTest(1, 0, 1*generator.GB, 1, printFile)
	defer clean(fileMap)
	if err != nil {
		t.Fatal(err)
	}

	// run command
	if err = executer.Execute(fileMap, "run"); err != nil {
		t.Fatal(err)
	}
	// check running ouput
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* migrate started", 1, printReg); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* has been finished", 1, printReg); err != nil {
		t.Fatal(err)
	}

	// check file equal
	if err := executer.CheckDirectroyEqual(fileMap); err != nil {
		t.Fatal(err)
	}
}

func TestManyFile(t *testing.T) {
	// env set
	fileMap, clean, err := generator.PrepareDefinedTest(10000, 0, 1, 1, printFile)
	defer clean(fileMap)
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
	fileMap, clean, err := generator.PrepareDefinedTest(1, 1, 3*generator.MB, 257, printFile)
	defer clean(fileMap)
	err = generator.CreateLocalDstDir(fileMap)
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
		"Start copying single object [A-Z0-9/]*/TESTFILE\\d+.camel", 257, true); err != nil {
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
	fileMap, clean, err := generator.PrepareDefinedTest(10, 3, 2*generator.MB, 3, printFile)
	defer clean(fileMap)
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
