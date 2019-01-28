package checker

import (
	"github.com/yunify/qscamel/tests/executer"
	"github.com/yunify/qscamel/tests/generater"
	"testing"
)

func TestDefaultRunCopy(t *testing.T) {

	fileMap, err := generater.CreateTestDefaultFile(
		"copy", "fs", "fs", nil, nil, true)
	//defer generater.CleanTestTempFile(fileMap)
	err = generater.CreateLocalSrcTestRandDirFile(
		fileMap, 4, 1, generater.MB*2, 2, true)
	err = generater.CreateLocalDstDir(fileMap)
	if err != nil {
		t.Fatal(err)
	}

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
