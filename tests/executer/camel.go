package executer

import (
	"fmt"
	"github.com/Netflix/go-expect"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"
)

func ExpectOutput(t *testing.T, expectString []string, fmap map[string]string) error{
	// check output if right
	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		return err
	}
	defer c.Close()
	_, file := path.Split(fmap["task"])
	sp := strings.Split(file, ".")
	cmd := exec.Command("qscamel", "run", sp[0], "-c", fmap["config"], "-t", fmap["task"], )
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	//done := make(chan int)

	switch len(expectString) {
	case 0 :
		go func() {
			c.ExpectEOF()

		}()


		if err = cmd.Start(); err != nil {
			return err
		}

		if err = cmd.Wait(); err != nil {
			return err
		}


	default:
		expt := make([]string, len(expectString))
		go func() {
			for i, ex := range expectString {
				fmt.Println(i, ex)
				expt[i], err = c.ExpectString(ex)

			}

		}()

		if err = cmd.Start(); err != nil {
			return err
		}
		if err = cmd.Wait(); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)

		for i, ex := range expt {
			if ex == "" {
				fmt.Println()
				return detected{fmt.Sprintf(
					"expect output: expect string '%s' not exsist\n",
					expectString[i])}
			} else {
				fmt.Println()
				fmt.Printf("'%s' got\n", expectString[i])
			}
		}
		fmt.Println()
	}
	return nil
}

type detected struct {
	fail string
}

func (e detected)Error() string{
	return fmt.Sprintf("%s", e.fail)
}
