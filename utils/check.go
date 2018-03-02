package utils

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

// CheckError will execute the func, handle it's panic and error
func CheckError(fn func() error) {
	defer func() {
		if x := recover(); x != nil {
			logrus.Fatalf("Caught panic: %v, Trace: %s", x, debug.Stack())
		}
	}()
	err := fn()
	if err != nil {
		logrus.Fatalf("Exited for error %v.", err)
	}
}
