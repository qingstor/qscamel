package executer

import (
	"fmt"
)

// ExpectOutput check the running result that the executing
// the command output to a file
func ExpectOutput() error {
	// check output if right
}

type detecter struct {
	fail string
}

func (e detecter) Error() string {
	return fmt.Sprintf("%s", e.fail)
}
