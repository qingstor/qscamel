package integration

import (
	"github.com/yunify/qscamel/tests/utils"
)

func TestDefaultRunCopy() {
	fileMap, clean := utils.PrepareDefaultTest()
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
	utils.CheckOutput(fileMap, "Task [a-z0-9]* has been finished", -1)
	utils.CheckDBNoObject(fileMap)

}

func TestDefaultDelete() {
	fileMap, clean := utils.PrepareDefaultTest()
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
	fileMap["delname"] = fileMap["name"]
	utils.Execute(fileMap, "delete")
	utils.CheckOutput(fileMap, "Task [a-z0-9]* has been deleted", -1)
	utils.CheckDBNoObject(fileMap)

}

func TestDefalutStatus() {
	fileMap, clean := utils.PrepareDefaultTest()
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
	utils.Execute(fileMap, "status")
	utils.CheckOutput(fileMap, "Show status started", 1)
	utils.CheckOutput(fileMap, "There are [0-9]* tasks totally", -1)
	utils.CheckDBNoObject(fileMap)
}

func TestDefaultClean() {
	fileMap, clean := utils.PrepareDefaultTest()
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
	utils.Execute(fileMap, "clean")
	utils.CheckOutput(fileMap, "Clean started", 1)
	utils.CheckOutput(fileMap, "Task [a-z0-9]* has been cleaned", -1)
	utils.CheckDBNoObject(fileMap)
}
