package s3

import (
	"context"
	"fmt"
	"strings"

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
	cp := utils.RebuildPath(c.Path, j.Key)

	fmt.Println("list cp: ", cp)

	marker := j.Marker

	// If ListObjectV2 is enabled, we should use ListObjectsV2 instead.
	if c.EnableListObjectsV2 {
		for {
			resp, err := c.client.ListObjectsV2(&s3.ListObjectsV2Input{
				Bucket:     aws.String(c.BucketName),
				Prefix:     aws.String(cp),
				MaxKeys:    aws.Int64(MaxKeys),
				StartAfter: aws.String(marker),
			})
			if err != nil {
				logrus.Errorf("List objects failed for %v.", err)
				return err
			}
			for _, v := range resp.Contents {
				if strings.HasSuffix(*v.Key, "/") {
					key := utils.GetRelativePathStrict(c.Path, *v.Key)
					_, err := c.client.HeadObject(&s3.HeadObjectInput{
						Bucket: aws.String(c.BucketName),
						Key:    aws.String(*v.Key),
					})
					if err == nil {
						so := &model.SingleObject{
							Key:          key,
							Size:         *v.Size,
							LastModified: v.LastModified.Unix(),
							MD5:          strings.Trim(*v.ETag, "\""),
							IsDir:        true,
						}
						fn(so)
					}
					continue
				}
				object := &model.SingleObject{
					Key:  utils.GetRelativePathStrict(c.Path, *v.Key),
					Size: *v.Size,
				}

				fn(object)
			}

			marker = aws.StringValue(resp.StartAfter)
			if !aws.BoolValue(resp.IsTruncated) {
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

	for {
		resp, err := c.client.ListObjects(&s3.ListObjectsInput{
			Bucket:  aws.String(c.BucketName),
			Prefix:  aws.String(cp),
			MaxKeys: aws.Int64(MaxKeys),
			Marker:  aws.String(marker),
		})
		if err != nil {
			logrus.Errorf("List objects failed for %v.", err)
			return err
		}
		for _, v := range resp.Contents {
			if strings.HasSuffix(*v.Key, "/") {
				key := utils.GetRelativePathStrict(c.Path, *v.Key)
				_, err := c.client.HeadObject(&s3.HeadObjectInput{
					Bucket: aws.String(c.BucketName),
					Key:    aws.String(*v.Key),
				})
				if err == nil {
					so := &model.SingleObject{
						Key:          key,
						Size:         *v.Size,
						LastModified: v.LastModified.Unix(),
						MD5:          strings.Trim(*v.ETag, "\""),
						IsDir:        true,
					}
					fn(so)
				}
				continue
			}
			object := &model.SingleObject{
				Key:  utils.GetRelativePathStrict(c.Path, *v.Key),
				Size: *v.Size,
			}

			fn(object)
		}

		// Some s3 compatible services may not return next marker, we should also try
		// the last key in contents.
		marker = aws.StringValue(resp.NextMarker)
		if marker == "" && len(resp.Contents) > 0 {
			marker = aws.StringValue(resp.Contents[len(resp.Contents)-1].Key)
		}

		if !aws.BoolValue(resp.IsTruncated) {
			marker = ""
		}

		if marker == "" {
			break
		}

		// Update task content.
		j.Marker = marker
		err = model.CreateObject(ctx, j)
		if err != nil {
			logrus.Errorf("Save task failed for %v.", err)
			return err
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
