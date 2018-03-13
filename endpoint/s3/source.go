package s3

import (
	"context"
	"io"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/yunify/qscamel/model"
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
func (c *Client) List(ctx context.Context, p string) (o []model.Object, err error) {
	o = []model.Object{}

	// Add "/" to list specific prefix.
	cp := path.Join(c.Path, p) + "/"
	// Trim left "/" to prevent object start with "/"
	cp = strings.TrimLeft(cp, "/")

	marker := ""
	first := true

	for marker != "" || first {
		resp, err := c.client.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:     aws.String(c.BucketName),
			Prefix:     aws.String(cp),
			MaxKeys:    aws.Int64(MaxKeys),
			Delimiter:  aws.String("/"),
			StartAfter: aws.String(marker),
		})
		if err != nil {
			return nil, err
		}
		for _, v := range resp.Contents {
			object := model.Object{
				Key:   path.Join(p, path.Base(*v.Key)),
				IsDir: false,
				Size:  *v.Size,
			}

			o = append(o, object)
		}
		for _, v := range resp.CommonPrefixes {
			object := model.Object{
				Key:   path.Join(p, path.Base(*v.Prefix)),
				IsDir: true,
				Size:  0,
			}

			o = append(o, object)
		}

		first = false
		marker = *resp.NextContinuationToken
		if !*resp.IsTruncated {
			marker = ""
		}
	}

	return
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string) (r io.ReadCloser, err error) {
	cp := path.Join(c.Path, p)
	// Trim left "/" to prevent object start with "/"
	cp = strings.TrimLeft(cp, "/")

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
