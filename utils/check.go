package utils

import (
	"os"
	"runtime/debug"
	"syscall"

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

// CheckExist return error if file not exist.
func CheckExist(path string) error {
	_, err := os.Stat(path)
	return err
}

// CheckWritable return error if the directory can not be write.
func CheckWritable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return false, err
		}
	}

	if info.Mode().Perm()&0200 == 0 {
		return false, &os.PathError{Op: "write", Path: path, Err: syscall.EPERM}

	}
	return true, nil
}
