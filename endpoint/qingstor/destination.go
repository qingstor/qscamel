package qingstor

import (
	"context"
	"io"
	"path"
	"strings"

	"github.com/pengsrc/go-shared/convert"
	"github.com/sirupsen/logrus"
	"github.com/yunify/qingstor-sdk-go/service"

	"github.com/yunify/qscamel/model"
)

// Fetchable implement destination.Fetchable
func (q *QingStor) Fetchable() bool {
	return true
}

// Writable implement destination.Writable
func (q *QingStor) Writable() bool {
	return true
}

// Write implement destination.Write
func (q *QingStor) Write(ctx context.Context, p string, r io.Reader) (err error) {
	t, err := model.GetTask(ctx)
	if err != nil {
		logrus.Panic(err)
	}

	cp := path.Join(t.Dst.Path, p)
	cp = strings.TrimLeft(cp, "/")
	if cp == "" {
		return
	}

	o, err := model.GetObject(ctx, p)
	if err != nil {
		logrus.Panic(err)
	}

	_, err = q.client.PutObject(cp, &service.PutObjectInput{
		Body:          r,
		ContentLength: convert.Int64(o.Size),
	})
	if err != nil {
		return
	}

	logrus.Debugf("QingStor wrote object %s.", cp)
	return
}

// Fetch implement destination.Fetch
func (q *QingStor) Fetch(ctx context.Context, p string) (err error) {
	return
}

// Dir implement destination.Dir
func (q *QingStor) Dir(ctx context.Context, p string) (err error) {
	t, err := model.GetTask(ctx)
	if err != nil {
		logrus.Panic(err)
	}

	cp := path.Join(t.Dst.Path, p)
	cp = strings.TrimLeft(cp, "/")
	if cp == "" {
		return
	}

	_, err = q.client.PutObject(cp, &service.PutObjectInput{
		ContentType: convert.String(DirectoryContentType),
	})
	if err != nil {
		return
	}

	logrus.Debugf("QingStor created dir %s.", cp)
	return
}
