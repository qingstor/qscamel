package qingstor

import (
	"context"
	"strings"

	"github.com/pengsrc/go-shared/convert"
	"github.com/qingstor/qingstor-sdk-go/v4/service"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.DirectoryObject, fn func(o model.Object)) (err error) {
	cp := utils.RebuildPath(c.Path, j.Key)

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
				key := utils.GetRelativePathStrict(c.Path, *v.Key)
				output, err := c.client.HeadObject(*v.Key, nil)
				if err == nil {
					so := &model.SingleObject{
						Key:          key,
						Size:         *v.Size,
						LastModified: int64(*v.Modified),
						MD5:          strings.Trim(*v.Etag, "\""),
						IsDir:        true,
					}
					so.QSMetadata = make(map[string]string)
					if c.UserDefineMeta && output.XQSMetaData != nil {
						so.QSMetadata = *output.XQSMetaData
					}
					if output.ContentType != nil {
						so.QSMetadata["ContentType"] = *output.ContentType
					}
					fn(so)
				}
				continue
			}

			key := utils.GetRelativePathStrict(c.Path, *v.Key)
			output, err := c.client.HeadObject(*v.Key, nil)
			if err == nil {
				object := &model.SingleObject{
					Key:          key,
					Size:         *v.Size,
					LastModified: int64(*v.Modified),
					MD5:          strings.Trim(*v.Etag, "\""),
				}
				object.QSMetadata = make(map[string]string)
				if c.UserDefineMeta && output.XQSMetaData != nil {
					object.QSMetadata = *output.XQSMetaData
				}
				if output.ContentType != nil {
					object.QSMetadata["ContentType"] = *output.ContentType
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
	cp := utils.RebuildPath(c.Path, p)

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
