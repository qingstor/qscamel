package fs

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
)

// Fetchable implement destination.Fetchable
func (f *Fs) Fetchable() bool {
	return false
}

// Writable implement destination.Writable
func (f *Fs) Writable() bool {
	return true
}

// Write implement destination.Write
func (f *Fs) Write(ctx context.Context, p string, r io.ReadCloser) (err error) {
	t, err := model.GetTask(ctx)
	if err != nil {
		logrus.Panic(err)
	}

	cp := path.Join(t.Dst.Path, p)

	file, err := os.Create(cp)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	if err != nil {
		return
	}

	logrus.Debugf("Fs wrote file %s.", cp)
	return
}

// Fetch implement destination.Fetch
func (f *Fs) Fetch(ctx context.Context, p string) (err error) {
	return
}

// Dir implement destination.Dir
func (f *Fs) Dir(ctx context.Context, p string) (err error) {
	t, err := model.GetTask(ctx)
	if err != nil {
		logrus.Panic(err)
	}

	cp := path.Join(t.Dst.Path, p)

	err = os.MkdirAll(cp, os.ModeDir|0777)
	if err != nil {
		return
	}

	logrus.Debugf("Fs created dir %s.", cp)
	return
}
