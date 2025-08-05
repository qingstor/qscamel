package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Deletable implement destination.Deletable
func (c *Client) Deletable() bool {
	return true
}

// Fetchable implement destination.Fetchable
func (c *Client) Fetchable() bool {
	return false
}

// Writable implement destination.Writable
func (c *Client) Writable() bool {
	return true
}

// Delete implement destination.Delete
func (c *Client) Delete(ctx context.Context, p string) (err error) {
	cp := utils.RebuildPath(c.Path, p)

	_, err = c.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(c.BucketName),
		Key:    aws.String(cp),
	})
	if err != nil {
		return
	}

	logrus.Debugf("s3 delete object %s.", cp)
	return
}

// Write implement destination.Write
func (c *Client) Write(ctx context.Context, p string, size int64, r io.Reader, isDir bool, _ map[string]string) (err error) {
	cp := utils.RebuildPath(c.Path, p)

	var input *s3.PutObjectInput
	if isDir {
		input = &s3.PutObjectInput{
			Bucket: aws.String(c.BucketName),
			Key:    aws.String(cp),
		}
	} else {
		input = &s3.PutObjectInput{
			Bucket: aws.String(c.BucketName),
			Key:    aws.String(cp),
			// wrap by limitReader to keep body consistent with size
			Body:          aws.ReadSeekCloser(io.LimitReader(r, size)),
			ContentLength: aws.Int64(size),
		}
	}
	_, err = c.client.PutObject(input)

	if err != nil {
		return
	}

	logrus.Debugf("s3 wrote object %s.", cp)
	return
}

// Fetch implement destination.Fetch
func (c *Client) Fetch(ctx context.Context, p, url string) (err error) {
	return constants.ErrEndpointFuncNotImplemented
}

// Partable implement destination.Partable
func (c *Client) Partable() bool {
	return true
}

// InitPart implement destination.InitPart
func (c *Client) InitPart(ctx context.Context, p string, size int64, _ map[string]string) (uploadID string, partSize int64, partNumbers int, err error) {
	cp := utils.RebuildPath(c.Path, p)

	resp, err := c.client.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String(c.BucketName),
		Key:    aws.String(cp),
	})

	if err != nil {
		return
	}

	uploadID = *resp.UploadId
	partSize, err = calculatePartSize(size)
	if err != nil {
		logrus.Errorf("Object %s is too large", p)
		return
	}

	partNumbers = int(size / partSize)
	if size%partSize != 0 {
		partNumbers++
	}
	return
}

// UploadPart implement destination.UploadPart
func (c *Client) UploadPart(ctx context.Context, o *model.PartialObject, r io.Reader) (err error) {
	cp := utils.RebuildPath(c.Path, o.Key)

	_, err = c.client.UploadPart(&s3.UploadPartInput{
		Bucket:        aws.String(c.BucketName),
		Key:           aws.String(cp),
		UploadId:      aws.String(o.UploadID),
		ContentLength: aws.Int64(o.Size),
		PartNumber:    aws.Int64(int64(o.PartNumber)),
		// wrap by limitReader to keep body consistent with size
		Body: aws.ReadSeekCloser(io.LimitReader(r, o.Size)),
	})
	if err != nil {
		return
	}

	// Trick: We need to check from current part number here.
	// if we check from -1, then complete will be skipped, because next will never be nil.
	next, err := model.NextPartialObject(ctx, o.Key, o.PartNumber)
	if err != nil {
		return
	}
	if next != nil {
		logrus.Debugf("s3 wrote partial object %s at %d.", o.Key, o.Offset)
		return nil
	}

	return nil
}

func (c *Client) CompleteParts(ctx context.Context, path string, uploadId string, totalNumber int) (err error) {
	cp := utils.RebuildPath(c.Path, path)

	parts := make([]*s3.CompletedPart, totalNumber)
	for i := 0; i < totalNumber; i++ {
		parts[i] = &s3.CompletedPart{
			PartNumber: aws.Int64(int64(i)),
		}
	}

	_, err = c.client.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(c.BucketName),
		Key:      aws.String(cp),
		UploadId: aws.String(uploadId),
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: parts,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AbortUploads(ctx context.Context, path string, uploadId string) (err error) {
	cp := utils.RebuildPath(c.Path, path)

	_, err = c.client.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
		Bucket:   aws.String(c.BucketName),
		Key:      aws.String(cp),
		UploadId: aws.String(uploadId),
	})
	if err != nil {
		return err
	}

	return
}
