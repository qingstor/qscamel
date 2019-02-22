package utils

import (
	"testing"

	"github.com/sirupsen/logrus"
)

type Tr struct {
	*testing.T
}

//var _ testing.TB = (*Tr)(nil)

func (t *Tr) Error(args ...interface{}) {
	logrus.Error(args...)
}
func (t *Tr) Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}
func (t *Tr)  Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}
func (t *Tr) Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
func (t *Tr) Log(args ...interface{}){
	logrus.Info(args...)
}
func (t *Tr) Logf(format string, args ...interface{}){
	logrus.Infof(format, args...)
}

// misimplements
func (t *Tr) Name() string {return ""}
func (t *Tr) Skip(args ...interface{}) {}
func (t *Tr) SkipNow() {}
func (t *Tr) Skipf(format string, args ...interface{}) {}
func (t *Tr) Skipped() bool { return false }
func (t *Tr) Helper() {}