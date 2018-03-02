package fs

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
)

// Reachable implement source.Reachable
func (f *Fs) Reachable() bool {
	return false
}

// Readable implement source.Readable
func (f *Fs) Readable() bool {
	return true
}

// List implement source.List
func (f *Fs) List(ctx context.Context, p string) (o []model.Object, err error) {
	task, err := model.GetTask(ctx)
	if err != nil {
		logrus.Panic(err)
	}

	cp := path.Join(task.Src.Path, p)

	fi, err := os.Open(cp)
	if err != nil {
		return nil, err
	}
	list, err := fi.Readdir(-1)
	fi.Close()

	o = make([]model.Object, len(list))
	for k, v := range list {
		o[k] = model.Object{
			Key:   path.Join(p, v.Name()),
			IsDir: v.IsDir(),
			Size:  v.Size(),
		}
	}

	return
}

// Read implement source.Read
func (f *Fs) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	task, err := model.GetTask(ctx)
	if err != nil {
		logrus.Panic(err)
	}

	cp := path.Join(task.Src.Path, p)

	r, err = os.Open(cp)
	if err != nil {
		logrus.Errorf("Fs open file %s failed for %s.", cp, err)
		return
	}
	return
}
