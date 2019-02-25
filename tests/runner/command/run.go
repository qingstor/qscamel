package command

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qscamel/tests/edge"
	"github.com/yunify/qscamel/tests/utils"
	"github.com/yunify/qscamel/tests/integration"
)

var RunCmd = &cobra.Command{
	Use:   "run [test name]",
	Short: "run a testing , like TestXXX",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		t := &utils.Tr{}
		switch args[0] {
		case "all":
			runAll(t)
		case "simple":
			runPart(t, int(SIMPLE))
		case "default":
			runPart(t, int(DEFAULT))
		case "dir":
			runPart(t, int(DIRECTORY))
		case "file":
			runPart(t, int(FILE))
		case "special":
			runPart(t, int(SPECIAL))
		case "endpoint":
			runPart(t, int(ENDPOINT))
		case "integration":
			runPart(t, int(SIMPLE), int(DEFAULT))
		case "edge":
			runPart(t, int(DIRECTORY), int(FILE),
				int(SPECIAL), int(ENDPOINT))
		default:
			runMatchCase(t, fmt.Sprintf("*%s*", args[0]))
		}
	},
}

var TestCase []map[string] func (t testing.TB)



func runAll(t testing.TB) {
	for i, _ := range TestCase {
		runPart(t, i)
	}
}

func runPart(t testing.TB, i... int) {
	for _, i := range i {
		for name, fn := range TestCase[i] {
			process(t, name, fn)
		}
	}

}

func runMatchCase(t testing.TB, Pattern string) {
	r, _ := regexp.Compile(Pattern)
	toTest := make(map[string] func(tb testing.TB))
	for _, set := range TestCase {
		for k, v := range set {
			if has := r.MatchString(k); has {
				toTest[k] = v
			}
		}
	}

	if len(toTest) == 0 {
		t.Log("No relevant tests")
		return
	}

	for k,v := range toTest {
		process(t, k, v)
	}
}

func process(t testing.TB, name string, fn func (t testing.TB)) {
	fmt.Printf("=== RUN   %s\n", name)
	s := time.Now()
	fn(t)
	fmt.Printf("--- PASS: %s (%.2fs)\n", name, time.Now().Sub(s).Seconds())
}

func init() {
	var TSIM,TDFT,TDIR,TFIL,TSPF,TEND map[string] func (t testing.TB)

	// Simple Test
	TSIM = map[string] func (t testing.TB) {
		"TestTaskRunCopy":  integration.TestTaskRunCopy,
		"TestTaskDelete": 	integration.TestTaskDelete,
		"TestTaskStatus": 	integration.TestTaskStatus,
		"TestTaskClean": 	integration.TestTaskClean,
	}

	// Default Test
	TDFT = map[string] func (t testing.TB) {
		"TestDefaultRunCopy": integration.TestDefaultRunCopy,
		"TestDefaultDelete":  integration.TestDefaultDelete,
		"TestDefalutStatus":  integration.TestDefalutStatus,
		"TestDefaultClean":   integration.TestDefaultClean,
	}

	// Directroy Test
	TDIR = map[string] func (t testing.TB) {
		"TestEmptyDirectory":	edge.TestEmptyDirectory,
		"TestOneDirectory": 	edge.TestOneDirectory,
		"TestDeepDirectory": 	edge.TestDeepDirectory,
		"TestManyDirectory": 	edge.TestManyDirectory,
	}

	// Normal File Test
	TFIL = map[string] func (t testing.TB) {
		"TestEmptyFile": 	  edge.TestEmptyFile,
		"TestBigFile": 		  edge.TestBigFile,
		"TestManyFile": 	  edge.TestManyFile,
		"TestDeepFile": 	  edge.TestDeepFile,
		"TestMutiDirAndFile": edge.TestMutiDirAndFile,
	}

	// Special File Test
	TSPF = map[string] func (t testing.TB) {
		"TestFileHole":    	   edge.TestFileHole,
		"TestDstSameFile": 	   edge.TestDstSameFile,
	}

	// Endpoint Test
	TEND = map[string] func (t testing.TB) {
		"TestFSInvalidDst":    edge.TestFSInvalidDst,
		"TestFSInvalidSrc":    edge.TestFSInvalidSrc,
	}

	TestCase = []map[string] func (t testing.TB){TSIM, TDFT, TDIR, TFIL, TSPF, TEND}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		DisableTimestamp: true,
	})
}

type TestingSet int

const (
	SIMPLE TestingSet = iota
	DEFAULT
	DIRECTORY
	FILE
	SPECIAL
	ENDPOINT
)
