package fs

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
)

// Deletable implement destination.Deletable
func (c *Client) Deletable() bool {
	return true
}

// Fetchable implement destination.Fetchable
func (c *Client) Fetchable() bool {
	return false
}

// Writable implement destination.Writable
func (c *Client) Writable() bool {
	return true
}

// Delete implement destination.Delete
func (c *Client) Delete(ctx context.Context, p string) (err error) {
	cp := filepath.Join(c.AbsPath, p)

	err = os.Remove(cp)
	if err != nil {
		return
	}

	logrus.Debugf("Fs delete file %s.", cp)
	return
}

// Write implement destination.Write
func (c *Client) Write(ctx context.Context, p string, _ int64, r io.Reader, _ bool, _ map[string]string) (err error) {
	cp, err := c.Encode(filepath.Join(c.AbsPath, p))
	if err != nil {
		return
	}

	_, err = os.Stat(filepath.Dir(cp))
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(cp), os.ModeDir|0777)
		if err != nil {
			return
		}
		logrus.Debugf("Fs created dir %s.", path.Dir(cp))
	}

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
func (c *Client) Fetch(ctx context.Context, p, url string) (err error) {
	return
}

// Partable implement destination.Partable
func (c *Client) Partable() bool {
	return false
}

// InitPart implement destination.InitPart
func (c *Client) InitPart(ctx context.Context, p string, size int64, _ map[string]string) (uploadID string, partSize int64, partNumbers int, err error) {
	return "", 0, 0, nil
}

// UploadPart implement destination.UploadPart
func (c *Client) UploadPart(ctx context.Context, o *model.PartialObject, r io.Reader) (err error) {
	return nil
}

func (c *Client) CompleteParts(ctx context.Context, path string, uploadId string, totalNumber int) (err error) {
	return nil
}

func (c *Client) AbortUploads(ctx context.Context, path string, uploadId string) (err error) {
	return nil
}
