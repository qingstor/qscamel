package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Reachable implement source.Reachable
func (c *Client) Reachable() bool {
	return true
}

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.DirectoryObject, fn func(o model.Object)) (err error) {
	cp := utils.Join(c.Path, j.Key) + "/"

	marker := j.Marker

	for {
		resp, err := c.client.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:     aws.String(c.BucketName),
			Prefix:     aws.String(cp),
			MaxKeys:    aws.Int64(MaxKeys),
			Delimiter:  aws.String("/"),
			StartAfter: aws.String(marker),
		})
		if err != nil {
			logrus.Errorf("List objects failed for %v.", err)
			return err
		}
		for _, v := range resp.Contents {
			object := &model.SingleObject{
				Key:  utils.Relative(*v.Key, c.Path),
				Size: *v.Size,
			}

			fn(object)
		}
		for _, v := range resp.CommonPrefixes {
			object := &model.DirectoryObject{
				Key: utils.Relative(*v.Prefix, c.Path),
			}

			fn(object)
		}

		marker = *resp.NextContinuationToken
		if !*resp.IsTruncated {
			marker = ""
		}

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
	return "", constants.ErrEndpointFuncNotImplemented
}

// Readable implement source.Readable
func (c *Client) Readable() bool {
	return false
}
