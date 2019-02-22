package edge

import (
	"testing"

	"github.com/yunify/qscamel/tests/utils"
)

func TestEmptyDirectory(t testing.TB) {
	fileMap, clean := utils.PrepareDefinedTest(t, 0, 0, utils.MB*2, 1)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
	utils.CheckOutputUnexpect(t, fileMap, "Job /DIR[0-9]* listed")
}

func TestOneDirectory(t testing.TB) {
	fileMap, clean := utils.PrepareDefinedTest(t, 0, 1, utils.MB*2, 2)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
}

func TestDeepDirectory(t testing.TB) {
	fileMap, clean := utils.PrepareDefinedTest(t, 0, 1, utils.MB*2, 257)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
}

func TestManyDirectory(t testing.TB) {
	fileMap, clean := utils.PrepareDefinedTest(t, 0, 3, 0, 3)
	defer clean(fileMap)
	utils.Execute(t, fileMap, "run")
}
