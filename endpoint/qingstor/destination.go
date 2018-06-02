package qingstor

import (
	"context"
	"io"

	"github.com/pengsrc/go-shared/convert"
	"github.com/sirupsen/logrus"
	"github.com/yunify/qingstor-sdk-go/service"

	"github.com/yunify/qscamel/utils"
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
func (c *Client) Write(ctx context.Context, p string, size int64, r io.ReadCloser) (err error) {
	cp := utils.Join(c.Path, p)

	if size <= c.MultipartBoundarySize {
		_, err = c.client.PutObject(cp, &service.PutObjectInput{
			Body:            r,
			ContentLength:   convert.Int64(size),
			XQSStorageClass: convert.String(c.StorageClass),
		})
	} else {
		err = c.uploader.Upload(r, cp)
	}
	if err != nil {
		return
	}

	logrus.Debugf("QingStor wrote object %s.", cp)
	return
}

// Fetch implement destination.Fetch
func (c *Client) Fetch(ctx context.Context, p, url string) (err error) {
	cp := utils.Join(c.Path, p)

	_, err = c.client.PutObject(cp, &service.PutObjectInput{
		XQSFetchSource: convert.String(url),
	})
	if err != nil {
		return
	}

	logrus.Debugf("QingStor fetched object %s.", cp)
	return
}
