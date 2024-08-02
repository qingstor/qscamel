package qingstor

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pengsrc/go-shared/convert"
	qsErrors "github.com/qingstor/qingstor-sdk-go/v4/request/errors"
	"github.com/qingstor/qingstor-sdk-go/v4/service"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Name implement base.Read
func (c *Client) Name(ctx context.Context) (name string) {
	return "qingstor:" + c.BucketName
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string, isDir bool) (r io.Reader, err error) {
	if isDir {
		return nil, nil
	}
	cp := utils.Join(c.Path, p)

	resp, err := c.client.GetObject(cp, nil)
	if err != nil {
		return
	}

	r = resp.Body
	return
}

// ReadRange implement source.ReadRange
func (c *Client) ReadRange(
	ctx context.Context, p string, offset, size int64,
) (r io.Reader, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.GetObject(cp, &service.GetObjectInput{
		Range: convert.String(fmt.Sprintf("bytes=%d-%d", offset, offset+size-1)),
	})
	if err != nil {
		return
	}

	r = resp.Body
	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string, isDir bool) (o *model.SingleObject, err error) {
	cp := utils.Join(c.Path, p)
	if isDir {
		cp += "/"
	}

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
	o = &model.SingleObject{
		Key:  p,
		Size: convert.Int64Value(resp.ContentLength),
		MD5:  strings.Trim(convert.StringValue(resp.ETag), "\""),
	}
	if resp.LastModified != nil {
		o.LastModified = (*resp.LastModified).Unix()
	}
	return
}
