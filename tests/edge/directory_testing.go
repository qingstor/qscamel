package edge

import (
	"github.com/yunify/qscamel/tests/utils"
)

func TestEmptyDirectory() {
	fileMap, clean := utils.PrepareDefinedTest(0, 0, utils.MB*2, 1)
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
	utils.CheckOutputUnexpect(fileMap, "Job /DIR[0-9]* listed")
}

func TestOneDirectory() {
	fileMap, clean := utils.PrepareDefinedTest(0, 1, utils.MB*2, 2)
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
}

func TestDeepDirectory() {
	fileMap, clean := utils.PrepareDefinedTest(0, 1, utils.MB*2, 257)
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
}

func TestManyDirectory() {
	fileMap, clean := utils.PrepareDefinedTest(0, 3, 0, 3)
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
}
