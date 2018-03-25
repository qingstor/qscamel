package qingstor

import (
	"context"
	"io"
	"strings"

	"github.com/pengsrc/go-shared/convert"
	"github.com/sirupsen/logrus"
	"github.com/yunify/qingstor-sdk-go/service"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Reachable implement source.Reachable
func (c *Client) Reachable() bool {
	return true
}

// Readable implement source.Readable
func (c *Client) Readable() bool {
	return true
}

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.Job, fn func(o *model.Object)) (err error) {
	om := make(map[string]struct{})

	cp := utils.Join(c.Path, j.Path) + "/"
	if cp == "/" {
		cp = ""
	}

	marker := j.Marker

	for {
		resp, err := c.client.ListObjects(&service.ListObjectsInput{
			Prefix:    convert.String(cp),
			Marker:    convert.String(marker),
			Limit:     convert.Int(MaxListObjectsLimit),
			Delimiter: convert.String("/"),
		})
		if err != nil {
			logrus.Errorf("List objects failed for %v.", err)
			return err
		}
		// Both "xxx/" and "xxx" with directory content type should be treated as directory.
		// And in order to prevent duplicate job, we need to use set to filter them.
		for _, v := range resp.Keys {
			object := &model.Object{
				Key:          utils.Relative(*v.Key, c.Path),
				IsDir:        *v.MimeType == DirectoryContentType,
				Size:         *v.Size,
				LastModified: int64(*v.Modified),
				MD5:          strings.Trim(*v.Etag, "\""),
			}

			if _, ok := om[object.Key]; !ok {
				fn(object)
				om[object.Key] = struct{}{}
			}
		}
		for _, v := range resp.CommonPrefixes {
			object := &model.Object{
				Key:   utils.Relative(*v, c.Path),
				IsDir: true,
				Size:  0,
			}

			if _, ok := om[object.Key]; !ok {
				fn(object)
				om[object.Key] = struct{}{}
			}
		}

		marker = *resp.NextMarker

		// Update task content.
		j.Marker = marker
		err = j.Save(ctx)
		if err != nil {
			logrus.Errorf("Save task failed for %v.", err)
			return err
		}

		if marker == "" {
			break
		}
	}

	return
}

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

// Reach implement source.Fetch
func (c *Client) Reach(ctx context.Context, p string) (url string, err error) {
	return
}
