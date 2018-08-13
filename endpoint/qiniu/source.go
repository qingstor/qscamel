package qiniu

import (
	"context"
	"time"

	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/x/rpc.v7"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.DirectoryObject, fn func(o model.Object)) (err error) {
	cp := utils.Join(c.Path, j.Key)
	if cp != "" {
		cp += "/"
	}

	marker := j.Marker

	for {
		entries, _, nextMarker, _, err := c.bucket.ListFiles(c.BucketName, cp, "", marker, MaxListFileLimit)
		if err != nil {
			if e, ok := err.(*rpc.ErrorInfo); ok {
				// If object not found, we just need to return a nil object.
				if e.Code == ErrorCodeInvalidMarker {
					marker = ""

					// Update task content.
					j.Marker = marker
					err = model.CreateObject(ctx, j)
					if err != nil {
						logrus.Errorf("Save task failed for %v.", err)
						return err
					}

					logrus.Warn("Qiniu's marker has been invalid, rescan for now.")
					continue
				}
			}

			logrus.Errorf("List files failed for %v.", err)
			return err
		}
		for _, v := range entries {
			object := &model.SingleObject{
				Key:  utils.Relative(v.Key, c.Path),
				Size: v.Fsize,
			}

			fn(object)
		}

		marker = nextMarker

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

	deadline := time.Now().Add(time.Hour).Unix()
	url = storage.MakePrivateURL(c.mac, c.Domain, cp, deadline)
	return
}

// Reachable implement source.Reachable
func (c *Client) Reachable() bool {
	return true
}
