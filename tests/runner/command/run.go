package command

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yunify/qscamel/tests/edge"
	"github.com/yunify/qscamel/tests/integration"
)

var RunCmd = &cobra.Command{
	Use:   "run [test name]",
	Short: "run a testing , like TestXXX",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "all":
			runAll()
		case "simple":
			runPart(int(SIMPLE))
		case "default":
			runPart(int(DEFAULT))
		case "dir":
			runPart(int(DIRECTORY))
		case "file":
			runPart(int(FILE))
		case "special":
			runPart(int(SPECIAL))
		case "endpoint":
			runPart(int(ENDPOINT))
		case "integration":
			runPart(int(SIMPLE), int(DEFAULT))
		case "edge":
			runPart(int(DIRECTORY), int(FILE),
				int(SPECIAL))
		default:
			runMatchCase(fmt.Sprintf("*%s*", args[0]))
		}
	},
}

func runAll() {
	for i, _ := range TestCase {
		runPart(i)
	}
}

func runPart(i... int) {
	for _, i := range i {
		for name, fn := range TestCase[i] {
			process(name, fn)
		}
	}

}

func runMatchCase(Pattern string) {
	r, _ := regexp.Compile(Pattern)
	toTest := make(map[string] func())
	for _, set := range TestCase {
		for k, v := range set {
			if has := r.MatchString(k); has {
				toTest[k] = v
			}
		}
	}

	if len(toTest) == 0 {
		log.Printf("No relevant tests")
		return
	}

	for k,v := range toTest {
		process(k, v)
	}
}

func process(name string, fn func ()) {
	fmt.Printf("=== RUN   %s\n", name)
	s := time.Now()
	fn()
	fmt.Printf("--- PASS: %s (%.2fs)\n", name, time.Now().Sub(s).Seconds())
}

var TestCase []map[string] func ()

func init() {
	var TSIM,TDFT,TDIR,TFIL,TSPF,TEND map[string] func ()

	// Simple Test
	TSIM = map[string] func () {
		"TestTaskRunCopy":  integration.TestTaskRunCopy,
		"TestTaskDelete": 	integration.TestTaskDelete,
		"TestTaskStatus": 	integration.TestTaskStatus,
		"TestTaskClean": 	integration.TestTaskClean,
	}

	// Default Test
	TDFT = map[string] func () {
		"TestDefaultRunCopy": integration.TestDefaultRunCopy,
		"TestDefaultDelete":  integration.TestDefaultDelete,
		"TestDefalutStatus":  integration.TestDefalutStatus,
		"TestDefaultClean":   integration.TestDefaultClean,
	}

	// Directroy Test
	TDIR = map[string] func () {
		"TestEmptyDirectory":	edge.TestEmptyDirectory,
		"TestOneDirectory": 	edge.TestOneDirectory,
		"TestDeepDirectory": 	edge.TestDeepDirectory,
		"TestManyDirectory": 	edge.TestManyDirectory,
	}

	// Normal File Test
	TFIL = map[string] func () {
		"TestEmptyFile": 	  edge.TestEmptyFile,
		"TestBigFile": 		  edge.TestBigFile,
		"TestManyFile": 	  edge.TestManyFile,
		"TestDeepFile": 	  edge.TestDeepFile,
		"TestMutiDirAndFile": edge.TestMutiDirAndFile,
	}

	// Special File Test
	TSPF = map[string] func () {
		"TestFileHole":    	   edge.TestFileHole,
		"TestDstSameFile": 	   edge.TestDstSameFile,
	}

	// Endpoint Test
	TEND = map[string] func () {
		"TestFSInvalidDst":    edge.TestFSInvalidDst,
		"TestFSInvalidSrc":    edge.TestFSInvalidSrc,
	}

	TestCase = []map[string] func (){TSIM, TDFT, TDIR, TFIL, TSPF, TEND}

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
