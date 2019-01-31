package edge

import (
	"testing"

	"github.com/yunify/qscamel/tests/utils/executer"
	"github.com/yunify/qscamel/tests/utils/generator"
)

func TestEmptyDirectory(t *testing.T) {
	// env set
	fileMap, clean, err := generator.PrepareDefinedTest(0, 0, generator.MB*2, 1, printFile)
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
		"Task task[0-9]* migrate started", 1, printReg); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* has been finished", 1, printReg); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutputUnexpect(fileMap, "Job /DIR[0-9]* listed", printReg); err != nil {
		t.Fatal(err)
	}

}

func TestOneDirectory(t *testing.T) {
	fileMap, clean, err := generator.PrepareDefinedTest(0, 1, generator.MB*2, 2, printFile)
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
		"Task task[0-9]* migrate started", 1, printReg); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutput(fileMap,
		"Task task[0-9]* has been finished", 1, printReg); err != nil {
		t.Fatal(err)
	}

}

func TestDeepDirectory(t *testing.T) {
	fileMap, clean, err := generator.PrepareDefinedTest(0, 1, generator.MB*2, 257, printFile)
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
	if err := executer.CheckOutput(fileMap,
		"Job (/DIR[0-9]*)* listed", 256, false); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutputUnexpect(fileMap,
		"Start copying single object [A-Z0-9]*/TESTFILE\\d+.camel", printReg); err != nil {
		t.Fatal(err)
	}
}

func TestManyDirectory(t *testing.T) {
	fileMap, clean, err := generator.PrepareDefinedTest(0, 3, 0, 3, printFile)
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
	if err := executer.CheckOutput(fileMap,
		"Job (/DIR[0-9]*)* listed", 36, false); err != nil {
		t.Fatal(err)
	}
	if err := executer.CheckOutputUnexpect(fileMap,
		"Start copying single object [A-Z0-9]*/TESTFILE\\d+.camel", printReg); err != nil {
		t.Fatal(err)
	}
}
