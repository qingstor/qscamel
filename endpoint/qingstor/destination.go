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
func (c *Client) Fetchable() bool {
	return true
}

// Writable implement destination.Writable
func (c *Client) Writable() bool {
	return true
}

// Write implement destination.Write
func (c *Client) Write(ctx context.Context, p string, r io.ReadCloser) (err error) {
	cp := path.Join(c.Path, p)
	cp = strings.TrimLeft(cp, "/")
	if cp == "" {
		return
	}

	o, err := model.GetObject(ctx, p)
	if err != nil {
		logrus.Panic(err)
	}

	_, err = c.client.PutObject(cp, &service.PutObjectInput{
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
func (c *Client) Fetch(ctx context.Context, p, url string) (err error) {
	cp := path.Join(c.Path, p)
	cp = strings.TrimLeft(cp, "/")
	if cp == "" {
		return
	}

	_, err = c.client.PutObject(cp, &service.PutObjectInput{
		XQSFetchSource: convert.String(url),
	})
	if err != nil {
		return
	}

	logrus.Debugf("QingStor fetched object %s.", cp)
	return
}

// Dir implement destination.Dir
func (c *Client) Dir(ctx context.Context, p string) (err error) {
	cp := path.Join(c.Path, p)
	cp = strings.TrimLeft(cp, "/")
	if cp == "" {
		return
	}

	_, err = c.client.PutObject(cp, &service.PutObjectInput{
		ContentType: convert.String(DirectoryContentType),
	})
	if err != nil {
		return
	}

	logrus.Debugf("QingStor created dir %s.", cp)
	return
}
