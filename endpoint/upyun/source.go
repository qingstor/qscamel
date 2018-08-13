package upyun

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/upyun/go-sdk/upyun"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.DirectoryObject, fn func(o model.Object)) (err error) {

	cp := utils.Join(c.Path, j.Key) + "/"

	oc := make(chan *upyun.FileInfo, 100)

	err = c.client.List(&upyun.GetObjectsConfig{
		Path:         cp,
		MaxListLevel: 1,
		ObjectsChan:  oc,
	})
	if err != nil {
		logrus.Errorf("List failed for %v.", err)
		return err
	}

	for v := range oc {
		if v.IsDir {
			o := &model.DirectoryObject{
				Key: utils.Relative(v.Name, c.Path),
			}

			fn(o)

			continue
		}
		o := &model.SingleObject{
			Key:  utils.Relative(v.Name, c.Path),
			Size: v.Size,
		}

		fn(o)
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
