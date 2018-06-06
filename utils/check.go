package utils

import (
	"os"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
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

// CheckClosedDB will check whether the error is a db closed err.
func CheckClosedDB(err error) {
	if err == leveldb.ErrClosed {
		logrus.Infof("Database has been closed, exit for now.")
		os.Exit(0)
	}
	logrus.Error(err)
}

// Recover will recover a goroutine from panic.
func Recover() {
	if x := recover(); x != nil {
		logrus.Fatalf("Caught panic: %v, Trace: %s", x, debug.Stack())
	}
}
