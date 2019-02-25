package edge

import (
	"os"
	"testing"

	"github.com/yunify/qscamel/tests/utils"
)

// directory isn't exist
func TestFSInvalidSrc(t testing.TB) {
	// env set
	fileMap, clean := utils.PrepareNormalTest(t)
	defer clean(fileMap)

	err := os.RemoveAll(fileMap["src"])
	if err != nil {
		t.Fatal(err)
	}

	// run command
	utils.Execute(t, fileMap, "run")
	// check output
	utils.CheckOutput(t, fileMap, "no such file or directory", 1)
}

func TestFSInvalidDst(t testing.TB) {
	// env set
	fileMap, clean := utils.PrepareNormalTest(t)
	defer clean(fileMap)

	//make destination not writable
	err := os.Chmod(fileMap["dst"], 0555)
	if err != nil {
		t.Fatal(err)
	}

	// run command
	utils.Execute(t, fileMap, "run")
	// check output
	utils.CheckOutput(t, fileMap, "operation not permitted", 1)
}