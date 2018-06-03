package qingstor

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pengsrc/go-shared/convert"
	qsErrors "github.com/yunify/qingstor-sdk-go/request/errors"
	"github.com/yunify/qingstor-sdk-go/service"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.GetObject(cp, nil)
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

	resp, err := c.client.GetObject(cp, &service.GetObjectInput{
		Range: convert.String(fmt.Sprintf("bytes=%d-%d", start, end)),
	})
	if err != nil {
		return
	}
	r := resp.Body
	defer r.Close()

	b = make([]byte, end-start+1)
	_, err = r.Read(b)
	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.Object, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.HeadObject(cp, nil)
	if err != nil {
		if e, ok := err.(*qsErrors.QingStorError); ok {
			// If object not found, we just need to return a nil object.
			if e.StatusCode == 404 {
				return nil, nil
			}
		}
		return
	}
	o = &model.Object{
		Key:   p,
		IsDir: convert.StringValue(resp.ContentType) == DirectoryContentType,
		Size:  convert.Int64Value(resp.ContentLength),
		MD5:   strings.Trim(convert.StringValue(resp.ETag), "\""),
	}
	if resp.LastModified != nil {
		o.LastModified = (*resp.LastModified).Unix()
	}
	return
}
