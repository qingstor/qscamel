package edge

import (
	"os/exec"
	"testing"

	"github.com/yunify/qscamel/tests/utils"
)

func TestFileHole(t testing.TB) {
	// env set
	fileMap := utils.CreateTestConfigFile(t, "copy", "fs", "fs", nil, nil)
	defer utils.CleanTestTempFile(fileMap)
	utils.CreateLocalSrcDir(t, fileMap)
	utils.CreateLocalDstDir(t, fileMap)
	utils.CreateHoleFile(t, fileMap["dir"], 2*utils.MB, 1*utils.MB, 1*utils.MB, 3)
	utils.CheckDirectroyEqual(t, fileMap)

}

func TestDstSameFile(t testing.TB) {
	// env set
	fileMap := utils.CreateTestConfigFile(t, "copy", "fs", "fs", nil, nil)
	defer utils.CleanTestTempFile(fileMap)
	utils.CreateLocalSrcDir(t, fileMap)
	utils.CreateLocalDstDir(t, fileMap)
	err := utils.CreateTestRandomFile(20, 2*utils.MB, fileMap["src"])
	cmd := exec.Command("cp", "-r", fileMap["src"], fileMap["dir"]+"/"+"dst")
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	// run command
	utils.Execute(t, fileMap, "run")
	// check file equal
	utils.CheckDirectroyEqual(t, fileMap)
}
