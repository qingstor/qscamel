package s3

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"io"

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
func (c *Client) List(ctx context.Context, j *model.Job, rc chan *model.Object) {
	defer close(rc)

	cp := utils.Join(c.Path, j.Path) + "/"

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
			rc <- nil
			return
		}
		for _, v := range resp.Contents {
			object := &model.Object{
				Key:   utils.Relative(*v.Key, c.Path),
				IsDir: false,
				Size:  *v.Size,
			}

			rc <- object
		}
		for _, v := range resp.CommonPrefixes {
			object := &model.Object{
				Key:   utils.Relative(*v.Prefix, c.Path),
				IsDir: true,
				Size:  0,
			}

			rc <- object
		}

		marker = *resp.NextContinuationToken
		if !*resp.IsTruncated {
			marker = ""
		}

		// Update task content.
		j.Marker = marker
		err = j.Save(ctx)
		if err != nil {
			logrus.Errorf("Save task failed for %v.", err)
			rc <- nil
			return
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

	resp, err := c.client.GetObject(&s3.GetObjectInput{
		Key:    aws.String(cp),
		Bucket: aws.String(c.BucketName),
	})
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
