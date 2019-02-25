package edge

import (
	"log"
	"os/exec"

	"github.com/yunify/qscamel/tests/utils"
)

func TestFileHole() {
	// env set
	fileMap := utils.CreateTestConfigFile("copy", "fs", "fs", nil, nil)
	defer utils.CleanTestTempFile(fileMap)
	utils.CreateLocalSrcDir(fileMap)
	utils.CreateLocalDstDir(fileMap)
	utils.CreateHoleFile(fileMap["dir"], 2*utils.MB, 1*utils.MB, 1*utils.MB, 3)
	utils.CheckDirectroyEqual(fileMap)

}

func TestDstSameFile() {
	// env set
	fileMap := utils.CreateTestConfigFile("copy", "fs", "fs", nil, nil)
	defer utils.CleanTestTempFile(fileMap)
	utils.CreateLocalSrcDir(fileMap)
	utils.CreateLocalDstDir(fileMap)
	err := utils.CreateTestRandomFile(20, 2*utils.MB, fileMap["src"])
	cmd := exec.Command("cp", "-r", fileMap["src"], fileMap["dir"]+"/"+"dst")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// run command
	utils.Execute(fileMap, "run")
	// check file equal
	utils.CheckDirectroyEqual(fileMap)
}
