package qingstor

import (
	"context"
	"strings"

	"github.com/pengsrc/go-shared/convert"
	"github.com/sirupsen/logrus"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.DirectoryObject, fn func(o model.Object)) (err error) {
	cp := utils.Join(c.Path, j.Key) + "/"
	if cp == "/" {
		cp = ""
	}

	marker := j.Marker

	for {
		resp, err := c.client.ListObjects(&service.ListObjectsInput{
			Prefix: convert.String(cp),
			Marker: convert.String(marker),
			Limit:  convert.Int(MaxListObjectsLimit),
		})
		if err != nil {
			logrus.Errorf("List objects failed for %v.", err)
			return err
		}

		for _, v := range resp.Keys {
			if strings.HasSuffix(*v.Key, "/") {
				key := utils.Relative(*v.Key, c.Path) + "/"
				output, err := c.client.HeadObject(*v.Key, nil)
				if err == nil {
					so := &model.SingleObject{
						Key:          key,
						Size:         *v.Size,
						LastModified: int64(*v.Modified),
						MD5:          strings.Trim(*v.Etag, "\""),
						IsDir:        true,
					}
					if c.UserDefineMeta {
						so.QSMetadata = output.XQSMetaData
					}
					fn(so)
				}
				continue
			}

			key := utils.Relative(*v.Key, c.Path)
			output, err := c.client.HeadObject(*v.Key, nil)
			if err == nil {
				object := &model.SingleObject{
					Key:          key,
					Size:         *v.Size,
					LastModified: int64(*v.Modified),
					MD5:          strings.Trim(*v.Etag, "\""),
				}
				if c.UserDefineMeta {
					object.QSMetadata = output.XQSMetaData
				}

				fn(object)
			}

		}

		marker = *resp.NextMarker

		// Update task content.
		j.Marker = marker
		err = model.CreateObject(ctx, j)
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

// Reach implement source.Fetch
func (c *Client) Reach(ctx context.Context, p string) (url string, err error) {
	cp := utils.Join(c.Path, p)

	r, _, err := c.client.GetObjectRequest(cp, nil)
	if err != nil {
		return
	}

	err = r.Build()
	if err != nil {
		return
	}

	err = r.SignQuery(3600)
	if err != nil {
		return
	}

	url = r.HTTPRequest.URL.String()
	return
}

// Reachable implement source.Reachable
func (c *Client) Reachable() bool {
	return true
}
