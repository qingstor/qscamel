package s3

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

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

// ReadAt implement source.ReadAt
func (c *Client) ReadAt(
	ctx context.Context, p string, start, end int64,
) (b []byte, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.GetObject(&s3.GetObjectInput{
		Key:    aws.String(cp),
		Bucket: aws.String(c.BucketName),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", start, end)),
	})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b = make([]byte, end-start+1)
	_, err = resp.Body.Read(b)
	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.Object, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(c.BucketName),
		Key:    aws.String(cp),
	})
	if err != nil {
		if e, ok := err.(awserr.Error); ok {
			if e.Code() == "NoSuchKey" {
				return nil, nil
			}
		}
		logrus.Errorf("Head object %s failed for %v.", p, err)
		return
	}
	o = &model.Object{
		Key:          p,
		IsDir:        strings.HasSuffix(p, "/"),
		Size:         *resp.ContentLength,
		MD5:          *resp.ETag,
		LastModified: (*resp.LastModified).Unix(),
	}
	return
}
