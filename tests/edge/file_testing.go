package edge

import (
	"github.com/yunify/qscamel/tests/utils"
)

func TestEmptyFile() {

	fileMap, clean := utils.PrepareDefinedTest(3, 0, 0, 1)
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
	// check file equal
	utils.CheckDirectroyEqual(fileMap)
}

// TODO: File Size should be adjusted appropriately
func TestBigFile() {
	fileMap, clean := utils.PrepareDefinedTest(1, 0, 1*utils.GB, 1)
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
	utils.CheckDirectroyEqual(fileMap)
}

func TestManyFile() {
	fileMap, clean := utils.PrepareDefinedTest(10000, 0, 1, 1)
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
	utils.CheckDirectroyEqual(fileMap)
}

func TestDeepFile() {
	fileMap, clean := utils.PrepareDefinedTest(1, 1, 3*utils.MB, 257)
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
	utils.CheckDirectroyEqual(fileMap)
}

func TestMutiDirAndFile() {
	fileMap, clean := utils.PrepareDefinedTest(10, 3, 2*utils.MB, 3,)
	defer clean(fileMap)
	utils.Execute(fileMap, "run")
	utils.CheckDirectroyEqual(fileMap)
}
