package edge

import (
	"os/exec"
	"testing"

	"github.com/yunify/qscamel/tests/utils/executer"
	"github.com/yunify/qscamel/tests/utils/generator"
)

func TestFileHole(t *testing.T) {
	// env set
	fileMap, err := generator.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, printFile)
	defer generator.CleanTestTempFile(fileMap)
	err = generator.CreateLocalSrcDir(fileMap)
	err = generator.CreateHoleFile((*fileMap)["src"], 2*generator.MB, 1*generator.MB, 1*generator.MB, 3)
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
		"Start copying single object [A-Z0-9/]*.hole", 3, printReg); err != nil {
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

func TestDstSameFile(t *testing.T) {
	// env set
	fileMap, err := generator.CreateTestConfigFile(
		"copy", "fs", "fs", nil, nil, printFile)
	defer generator.CleanTestTempFile(fileMap)
	err = generator.CreateLocalSrcDir(fileMap)
	err = generator.CreateTestRandomFile(20, 2*generator.MB, (*fileMap)["src"])
	cmd := exec.Command("cp", "-r", (*fileMap)["src"], (*fileMap)["dir"]+"/"+"dst")
	cmd.Run()
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
		"Start copying single object [A-Z0-9]*/TESTFILE\\d+.camel", 20, printReg); err != nil {
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
