package executer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// Execute base on task directory, executing the command
// on different platform, and the output will be redirected
// to a OXXXX.output
func Execute(fmap *map[string]string, comm string) error {

	// generate corrisponding argument to qscamel
	arg := strings.Join([]string{"-c", (*fmap)["config"]}, " ")
	if comm == "run" {
		arg = strings.Join([]string{"run", (*fmap)["name"], "-t", (*fmap)["task"], arg}, " ")
	}

	var c *exec.Cmd
	switch runtime.GOOS {
	default:
		c = cmdOnUnix("qscamel", strings.Split(arg, " ")...)
	}

	// set output file
	out, err := ioutil.TempFile((*fmap)["dir"], "O*.output")
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
	if err != nil {
		return err
	} else if len(*sl) != n {
		return detecter{fmt.Sprintf(
			"`%s` ouccurs %d times, expect %d",
			expectPattern, len(*sl), n)}
	}

	if p {
		for _, s := range *sl {
			fmt.Printf("'%s' got\n", s)
		}
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
