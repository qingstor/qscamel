package fs

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

// Fetchable implement destination.Fetchable
func (c *Client) Fetchable() bool {
	return false
}

// Writable implement destination.Writable
func (c *Client) Writable() bool {
	return true
}

// Write implement destination.Write
func (c *Client) Write(ctx context.Context, p string, r io.ReadCloser) (err error) {
	cp := path.Join(c.Path, p)

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
func (c *Client) Fetch(ctx context.Context, p string) (err error) {
	return
}

// Dir implement destination.Dir
func (c *Client) Dir(ctx context.Context, p string) (err error) {
	cp := path.Join(c.Path, p)

	err = os.MkdirAll(cp, os.ModeDir|0777)
	if err != nil {
		return
	}

	logrus.Debugf("Fs created dir %s.", cp)
	return
}
