package integration

import (
	"github.com/yunify/qscamel/tests/utils"
)

func TestTaskRunCopy() {
	// env set
	fileMap, clean := utils.PrepareNormalTest()
	defer clean(fileMap)

	// run command
	utils.Execute(fileMap, "run")

	// check running ouput
	utils.CheckOutput(fileMap, "Task [a-z0-9]* has been finished", 1)
	utils.CheckDBNoObject(fileMap)
}

func TestTaskDelete() {
	fileMap, clean := utils.PrepareNormalTest()
	defer clean(fileMap)
	utils.Execute(fileMap, "run")

	// delete command
	fileMap["delname"] = fileMap["name"]
	utils.Execute(fileMap, "delete")

	// check delete output
	utils.CheckOutput(fileMap, "Task [a-z0-9]* has been deleted", 1)
	utils.CheckDBNoObject(fileMap)
}

func TestTaskStatus() {
	fileMap, clean := utils.PrepareNormalTest()
	defer clean(fileMap)
	utils.Execute(fileMap, "run")

	// status command
	utils.Execute(fileMap, "status")

	// check status output
	utils.CheckOutput(fileMap, "Show status started", 1)
	utils.CheckOutput(fileMap, "There are 1 tasks totally", 1)
	utils.CheckDBNoObject(fileMap)
}

func TestTaskClean() {
	fileMap, clean := utils.PrepareNormalTest()
	defer clean(fileMap)
	utils.Execute(fileMap, "run")

	// clean command
	utils.Execute(fileMap, "clean")

	// check clean output
	utils.CheckOutput(fileMap, "Clean started", 1)
	utils.CheckOutput(fileMap, "Task [a-z0-9]* has been cleaned", 1)
	utils.CheckDBNoObject(fileMap)
}
