package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Name implement base.Read
func (c *Client) Name(ctx context.Context) (name string) {
	return "s3:" + c.BucketName
}

// Read implement source.Read
func (c *Client) Read(ctx context.Context, p string, isDir bool) (r io.Reader, err error) {
	if isDir {
		return nil, nil
	}
	cp := utils.RebuildPath(c.Path, p)
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

// ReadRange implement source.ReadRange
func (c *Client) ReadRange(
	ctx context.Context, p string, offset, size int64,
) (r io.Reader, err error) {
	cp := utils.RebuildPath(c.Path, p)

	resp, err := c.client.GetObject(&s3.GetObjectInput{
		Key:    aws.String(cp),
		Bucket: aws.String(c.BucketName),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", offset, offset+size-1)),
	})
	if err != nil {
		return
	}

	return resp.Body, nil
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string, isDir bool) (o *model.SingleObject, err error) {
	cp := utils.RebuildPath(c.Path, p)

	resp, err := c.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(c.BucketName),
		Key:    aws.String(cp),
	})
	if err != nil {
		switch e := err.(type) {
		case awserr.RequestFailure:
			if e.StatusCode() == 404 {
				return nil, nil
			}
		case awserr.Error:
			if e.Code() == "NoSuchKey" {
				return nil, nil
			}
		}
		logrus.Errorf("Head object %s failed for %v.", p, err)
		return
	}
	o = &model.SingleObject{
		Key:          p,
		Size:         *resp.ContentLength,
		MD5:          *resp.ETag,
		LastModified: (*resp.LastModified).Unix(),
	}
	return
}
