package edge

import (
	"testing"

	"github.com/yunify/qscamel/tests/utils"
)

func TestEmptyFile(t testing.TB) {

	fileMap, clean := utils.PrepareDefinedTest(t, 3, 0, 0, 1)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
	// check file equal
	utils.CheckDirectroyEqual(t, fileMap)
}

// TODO: File Size should be adjusted appropriately
func TestBigFile(t testing.TB) {
	fileMap, clean := utils.PrepareDefinedTest(t, 1, 0, 1*utils.GB, 1)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
	utils.CheckDirectroyEqual(t, fileMap)
}

func TestManyFile(t testing.TB) {
	fileMap, clean := utils.PrepareDefinedTest(t, 10000, 0, 1, 1)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
	utils.CheckDirectroyEqual(t, fileMap)
}

func TestDeepFile(t testing.TB) {
	fileMap, clean := utils.PrepareDefinedTest(t, 1, 1, 3*utils.MB, 257)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
	utils.CheckDirectroyEqual(t, fileMap)
}

func TestMutiDirAndFile(t testing.TB) {
	fileMap, clean := utils.PrepareDefinedTest(t, 10, 3, 2*utils.MB, 3,)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
	utils.CheckDirectroyEqual(t, fileMap)
}
