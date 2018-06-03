package aliyun

import (
	"context"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.Job, fn func(o *model.Object)) (err error) {
	cp := utils.Join(c.Path, j.Path) + "/"

	marker := j.Marker

	for {
		resp, err := c.client.ListObjects(
			oss.Delimiter("/"),
			oss.Marker(marker),
			oss.MaxKeys(MaxKeys),
			oss.Prefix(cp),
		)
		if err != nil {
			logrus.Errorf("List objects failed for %v.", err)
			return err
		}
		for _, v := range resp.Objects {
			object := &model.Object{
				Key:   utils.Relative(v.Key, c.Path),
				IsDir: false,
				Size:  v.Size,
			}

			fn(object)
		}
		for _, v := range resp.CommonPrefixes {
			object := &model.Object{
				Key:   utils.Relative(v, c.Path),
				IsDir: true,
				Size:  0,
			}

			fn(object)
		}

		marker = resp.NextMarker

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

// Reach implement source.Fetch
func (c *Client) Reach(ctx context.Context, p string) (url string, err error) {
	return "", constants.ErrEndpointFuncNotImplemented
}

// Reachable implement source.Reachable
func (c *Client) Reachable() bool {
	return false
}
