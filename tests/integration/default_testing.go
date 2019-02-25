package integration

import (
	"testing"

	"github.com/yunify/qscamel/tests/utils"
)

func TestDefaultRunCopy(t testing.TB) {
	fileMap, clean := utils.PrepareDefaultTest(t)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
	utils.CheckOutput(t, fileMap, "Task [a-z0-9]* has been finished", 1)
	utils.CheckDBNoObject(t, fileMap)

}

func TestDefaultDelete(t testing.TB) {
	fileMap, clean := utils.PrepareDefaultTest(t)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
	fileMap["delname"] = fileMap["name"]
	utils.Execute(t, fileMap, "delete")
	utils.CheckOutput(t, fileMap, "Task [a-z0-9]* has been deleted", 1)
	utils.CheckDBNoObject(t, fileMap)

}

func TestDefalutStatus(t testing.TB) {
	fileMap, clean := utils.PrepareDefaultTest(t)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
	utils.Execute(t, fileMap, "status")
	utils.CheckOutput(t, fileMap, "Show status started", 1)
	utils.CheckOutput(t, fileMap, "There are [0-9]* tasks totally", 1)
	utils.CheckDBNoObject(t, fileMap)
}

func TestDefaultClean(t testing.TB) {
	fileMap, clean := utils.PrepareDefaultTest(t)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
	utils.Execute(t, fileMap, "clean")
	utils.CheckOutput(t, fileMap, "Clean started", 1)
	utils.CheckOutput(t, fileMap, "Task [a-z0-9]* has been cleaned", 1)
	utils.CheckDBNoObject(t, fileMap)
}
