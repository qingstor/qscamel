package executer

import (
	"fmt"
	"github.com/yunify/qscamel/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"qiniupkg.com/x/log.v7"
	"regexp"
	"runtime"
	"strings"
)

// Execute base on task directory, executing the command
// on different platform, and the output will be redirected
// to a 'comm'+XXXX.output
func Execute(fmap *map[string]string, comm string) error {

	var arg string
	// generate corrisponding argument to qscamel
	if _, has := (*fmap)["config"]; has == true {
		arg = strings.Join([]string{"-c", (*fmap)["config"]}, " ")
	}
	args := strings.Split(arg, " ")

	switch comm {
	case "run":
		arg = strings.Join([]string{"run", (*fmap)["name"], "-t", (*fmap)["task"], arg}, " ")
	case "delete":
		arg = strings.Join([]string{comm, (*fmap)["delname"], arg}, " ")
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
	out, err := ioutil.TempFile((*fmap)["dir"], comm+"*.output")
	if err != nil {
		return err
	}
	defer out.Close()

	(*fmap)["output"] = out.Name()
	c.Stdout = out
	c.Stderr = out

	// run command
	if err = c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

// CheckOutput will check the output file after executing a command
// and return error if the expect count 'n' is not equal to the count
// of satisfied string.
func CheckOutput(fmap *map[string]string, expectPattern string, n int, p bool) error {
	out, err := os.Open((*fmap)["output"])
	if err != nil {
		return err
	}
	defer out.Close()

	// check out put
	stm, err := ioutil.ReadAll(out)
	if err != nil {
		return err
	}
	sl, err := ExpectOutput(&stm, expectPattern)
	fmt.Printf("re: %s ... (expect: %d/got: %d)\n", expectPattern[:5], n, len(*sl))
	if err != nil {
		return err
	} else if len(*sl) != n {
		return detecter{fmt.Sprintf("not satisfied")}
	}
	if p {
		for _, s := range *sl {
			log.Infof("%s\n", s)
		}
	}
	return nil
}

// CheckOutputUnexpect will check the output file after executing a command
// and return error if the unexpected string has occurrences
func CheckOutputUnexpect(fmap *map[string]string, expectPattern string, p bool) error {
	out, err := os.Open((*fmap)["output"])
	if err != nil {
		return err
	}
	defer out.Close()

	// check out put
	stm, err := ioutil.ReadAll(out)
	if err != nil {
		return err
	}
	sl, err := ExpectOutput(&stm, expectPattern)
	if err != nil {
		return err
	} else if len(*sl) != 0 {

		if p {
			for _, s := range *sl {
				log.Infof("Unexpected string '%s'\n", s)
			}
		}
		return detecter{fmt.Sprintf("Unexpected string")}
	}

	return nil
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

// CheckDirectroyEqual check two dirctory if is equal
func CheckDirectroyEqual(fmap *map[string]string) error {
	eq, err := utils.CompareLocalDirectoryMD5((*fmap)["dir"]+"/src", (*fmap)["dir"]+"/dst")
	if err != nil {
		return err
	}
	if eq == false {
		return detecter{"check directory is finished: not equal"}
	}
	log.Info("check directory is finished: equal")
	return nil
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
