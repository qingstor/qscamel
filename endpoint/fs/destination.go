package fs

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/utils"
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
	cp := "/" + utils.Join(c.Path, p)

	_, err = os.Stat(path.Dir(cp))
	if os.IsNotExist(err) {
		err = os.MkdirAll(path.Dir(cp), os.ModeDir|0777)
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
