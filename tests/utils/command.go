package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// Execute base on task directory, executing the command
// on different platform, and the output will be redirected
// to a 'comm'+XXXX.output
func Execute(fmap map[string]string, comm string) {

	var arg string
	// generate corrisponding argument to qscamel
	if _, has := fmap["config"]; has == true {
		arg = strings.Join([]string{"-c", fmap["config"]}, " ")
	}
	args := strings.Split(arg, " ")

	switch comm {
	case "run":
		arg = strings.Join([]string{"run", fmap["name"], "-t", fmap["task"], arg}, " ")
	case "delete":
		arg = strings.Join([]string{comm, fmap["delname"], arg}, " ")
	default:
		arg = strings.Join([]string{comm, arg}, " ")
	}

	// remove ""
	args = strings.Split(arg, " ")
	if args[len(args)-1] == "" {
		args = args[:len(args)-1]
	}

	var c *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// TODO
	default:
		c = cmdOnUnix("qscamel", args...)
	}

	// set output file
	out, err := ioutil.TempFile(fmap["dir"], comm+"*.output")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	fmap["output"] = out.Name()
	c.Stdout = out
	c.Stderr = out

	// run command
	if err = c.Run(); err != nil {
		log.Fatal(err)
	}

}

// CheckOutput will check the output file after executing a command
// and fatal if the expect count 'n' is not equal to the count
// of satisfied string.
func CheckOutput(fmap map[string]string, expectPattern string, n int) {
	out, err := os.Open(fmap["output"])
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// check out put
	stm, err := ioutil.ReadAll(out)
	if err != nil {
		log.Fatal(err)
	}
	sl, err := ExpectOutput(&stm, expectPattern)
	log.Printf("regexp: %s ... (expect: %d/got: %d)\n", expectPattern[:5], n, len(*sl))
	if err != nil {
		log.Fatal(err)
	} else if len(*sl) != n {
		log.Fatal(detecter{fmt.Sprintf("not satisfied '%s'", expectPattern)})
	}

}

// CheckOutputUnexpect will check the output file after executing a command
// and fatal if the unexpected string has occurrences
func CheckOutputUnexpect(fmap map[string]string, expectPattern string) {
	out, err := os.Open(fmap["output"])
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// check out put
	stm, err := ioutil.ReadAll(out)
	if err != nil {
		log.Fatal(err)
	}
	sl, err := ExpectOutput(&stm, expectPattern)
	if err != nil {
		log.Fatal(err)
	} else if len(*sl) != 0 {
		for _, s := range *sl {
			log.Printf("Unexpected string '%s'\n", s)
		}

		log.Fatal(detecter{fmt.Sprintf("Unexpected string")})
	}

}

// ExpectOutput use the regular expession to check the byte stream
// return all strings satisfy the regex pattern
func ExpectOutput(content *[]byte, regex string) (*[]string, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	sl := re.FindAllString(string(*content), -1)
	return &sl, nil
}

func cmdOnWin(comm string, arg ...string) *exec.Cmd {
	// Something TODO
	return exec.Command(comm, arg...)
}

func cmdOnUnix(comm string, arg ...string) *exec.Cmd {
	return exec.Command(comm, arg...)
}

type detecter struct {
	fail string
}

func (e detecter) Error() string {
	return fmt.Sprintf("%s", e.fail)
}
