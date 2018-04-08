package fs

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Client is the struct for POSIX file system endpoint.
type Client struct {
	Path string
}

// New will create a Fs.
func New(ctx context.Context, et uint8) (c *Client, err error) {
	t, err := model.GetTask(ctx)
	if err != nil {
		return
	}

	e := t.Src
	if et == constants.DestinationEndpoint {
		e = t.Dst
	}

	c = &Client{}

	// Set prefix.
	c.Path = e.Path

	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.Object, err error) {
	cp := "/" + utils.Join(c.Path, p)

	fi, err := os.Stat(cp)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		logrus.Errorf("Stat %s failed for %v.", p, err)
		return
	}
	// We will not calculate md5 while stating object.
	o = &model.Object{
		Key:          p,
		IsDir:        fi.IsDir(),
		Size:         fi.Size(),
		LastModified: fi.ModTime().Unix(),
	}
	return
}

// MD5 implement source.MD5 and destination.MD5
func (c *Client) MD5(ctx context.Context, p string) (b string, err error) {
	r, err := c.Read(ctx, p)
	if err != nil {
		return
	}
	defer r.Close()

	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	sum := h.Sum(nil)

	return hex.EncodeToString(sum[:]), nil
}

// MD5able implement source MD5able and destination MD5able.
func (c *Client) MD5able() bool {
	return true
}
