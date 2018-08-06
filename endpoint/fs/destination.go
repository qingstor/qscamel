package fs

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
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
func (c *Client) Write(ctx context.Context, p string, _ int64, r io.ReadCloser) (err error) {
	cp := filepath.Join(c.AbsPath, p)

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
