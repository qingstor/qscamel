package qiniu

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/x/rpc.v7"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Name implement base.Read
func (c *Client) Name(ctx context.Context) (name string) {
	return "qiniu:" + c.BucketName
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := utils.Join(c.Path, p)

	deadline := time.Now().Add(time.Hour).Unix()
	url := storage.MakePrivateURL(c.mac, c.Domain, cp, deadline)

	resp, err := c.client.Get(url)
	if err != nil {
		return
	}

	r = resp.Body
	return
}

// ReadAt implement source.ReadAt
func (c *Client) ReadAt(
	ctx context.Context, p string, start, end int64,
) (b []byte, err error) {
	cp := utils.Join(c.Path, p)

	deadline := time.Now().Add(time.Hour).Unix()
	url := storage.MakePrivateURL(c.mac, c.Domain, cp, deadline)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b = make([]byte, end-start+1)
	_, err = resp.Body.Read(b)
	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.Object, err error) {
	cp := utils.Join(c.Path, p)

	fi, err := c.bucket.Stat(c.BucketName, cp)
	if err != nil {
		if e, ok := err.(*rpc.ErrorInfo); ok {
			// If object not found, we just need to return a nil object.
			if e.Code == ErrorCodeNotFound {
				return nil, nil
			}
		}
		logrus.Errorf("Stat failed for %v.", err)
		return
	}
	// qiniu use their own hash algorithm instead of md5, so we can't support it.
	o = &model.Object{
		Key:          p,
		IsDir:        strings.HasSuffix(p, "/"),
		Size:         fi.Fsize,
		LastModified: fi.PutTime,
	}
	return
}
