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
		case "dir":
			runPart(t, 0)
		case "file":
			runPart(t, 1)
		case "special":
			runPart(t, 2)
		case "endpoint":
			runPart(t, 3)
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

func runPart(t testing.TB, i int) {
	for name, fn := range TestCase[i] {
		process(t, name, fn)
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
	var TDIR, TFIL, TSPF, TEND map[string] func (t testing.TB)

	TDIR = map[string] func (t testing.TB) {
		"TestEmptyDirectory":	edge.TestEmptyDirectory,
		"TestOneDirectory": 	edge.TestOneDirectory,
		"TestDeepDirectory": 	edge.TestDeepDirectory,
		"TestManyDirectory": 	edge.TestManyDirectory,
	}
	TFIL = map[string] func (t testing.TB) {
		"TestEmptyFile": 	  edge.TestEmptyFile,
		"TestBigFile": 		  edge.TestBigFile,
		"TestManyFile": 	  edge.TestManyFile,
		"TestDeepFile": 	  edge.TestDeepFile,
		"TestMutiDirAndFile": edge.TestMutiDirAndFile,
	}
	TSPF = map[string] func (t testing.TB) {
		"TestFileHole":    edge.TestFileHole,
		"TestDstSameFile": edge.TestDstSameFile,
	}
	TEND = map[string] func (t testing.TB) {
		"TestFSInvalidDst":    edge.TestFSInvalidDst,
		"TestFSInvalidSrc": edge.TestFSInvalidSrc,
	}
	TestCase = []map[string] func (t testing.TB){TDIR, TFIL, TSPF, TEND}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		DisableTimestamp: true,
	})
}