package integration

import (
	"testing"

	"github.com/yunify/qscamel/tests/utils"
)

func TestTaskRunCopy(t testing.TB) {
	// env set
	fileMap, clean := utils.PrepareNormalTest(t)
	defer clean(fileMap)

	// run command
	utils.Execute(t, fileMap, "run")

	// check running ouput
	utils.CheckOutput(t, fileMap, "Task [a-z0-9]* has been finished", 1)
	utils.CheckDBNoObject(t, fileMap)
}

func TestTaskDelete(t testing.TB) {
	fileMap, clean := utils.PrepareNormalTest(t)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")

	// delete command
	fileMap["delname"] = fileMap["name"]
	utils.Execute(t, fileMap, "delete")

	// check delete output
	utils.CheckOutput(t, fileMap, "Task [a-z0-9]* has been deleted", 1)
	utils.CheckDBNoObject(t, fileMap)
}

func TestTaskStatus(t testing.TB) {
	fileMap, clean := utils.PrepareNormalTest(t)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")

	// status command
	utils.Execute(t, fileMap, "status")

	// check status output
	utils.CheckOutput(t, fileMap, "Show status started", 1)
	utils.CheckOutput(t, fileMap, "There are 1 tasks totally", 1)
	utils.CheckDBNoObject(t, fileMap)
}

func TestTaskClean(t testing.TB) {
	fileMap, clean := utils.PrepareNormalTest(t)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")

	// clean command
	utils.Execute(t, fileMap, "clean")

	// check clean output
	utils.CheckOutput(t, fileMap, "Clean started", 1)
	utils.CheckOutput(t, fileMap, "Task [a-z0-9]* has been cleaned", 1)
	utils.CheckDBNoObject(t, fileMap)
}
