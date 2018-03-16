package fs

import (
	"context"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
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
	cp := path.Join(c.Path, p)

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
