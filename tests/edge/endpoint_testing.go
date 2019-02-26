package edge

import (
	"log"
	"os"

	"github.com/yunify/qscamel/tests/utils"
)


func TestFSInvalidSrc() {
	// env set
	fileMap, clean := utils.PrepareNormalTest()
	defer clean(fileMap)

	// source directory isn't exist
	err := os.RemoveAll(fileMap["src"])
	if err != nil {
		log.Fatal(err)
	}

	// run command
	utils.Execute(fileMap, "run")
	// check output
	utils.CheckOutput(fileMap, "no such file or directory", 1)
	utils.CheckDBEmpty(fileMap)
}

func TestFSInvalidDst() {
	// env set
	fileMap, clean := utils.PrepareNormalTest()
	defer clean(fileMap)

	//make destination not writable
	err := os.Chmod(fileMap["dst"], 0555)
	if err != nil {
		log.Fatal(err)
	}

	// run command
	utils.Execute(fileMap, "run")
	// check output
	utils.CheckOutput(fileMap, "operation not permitted", 1)
	utils.CheckDBEmpty(fileMap)
}
