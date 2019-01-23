package checker

import (
	"github.com/yunify/qscamel/tests/dbtester"
	"github.com/yunify/qscamel/tests/executer"
	"github.com/yunify/qscamel/tests/generater"
	"testing"
)


func TestTaskRun(t *testing.T) {
	// env
	fileMap := generater.CreateTestConfigFile(t)
	generater.CreateTestRandDirFile(t, fileMap, 10, 3,
		generater.MB * 2, 3, true)

	// check
	if err := dbtester.CheckDBEmpty(fileMap); err != nil {
		generater.Fatal(t, fileMap, err)
	}
	if err := executer.ExpectOutput(t, []string{"Task", "has been finished. " }, fileMap); err != nil {
		generater.Fatal(t, fileMap, err)
	}

	// clean
	generater.CleanTestTempFile(t, fileMap)
}




