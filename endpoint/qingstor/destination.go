package qingstor

import (
	"context"
	"io"
	"strings"

	"github.com/pengsrc/go-shared/convert"
	"github.com/sirupsen/logrus"
	"github.com/yunify/qingstor-sdk-go/v3/service"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Deletable implement destination.Deletable
func (c *Client) Deletable() bool {
	return true
}

// Fetchable implement destination.Fetchable
func (c *Client) Fetchable() bool {
	return true
}

// Writable implement destination.Writable
func (c *Client) Writable() bool {
	return true
}

// Delete implement destination.Delete
func (c *Client) Delete(ctx context.Context, p string) (err error) {
	cp := utils.Join(c.Path, p)

	_, err = c.client.DeleteObject(cp)
	if err != nil {
		return
	}

	logrus.Debugf("QingStor delete object %s.", cp)
	return
}

// Write implement destination.Write
func (c *Client) Write(ctx context.Context, p string, size int64, r io.Reader, isDir bool, meta map[string]string) (err error) {
	cp := utils.Join(c.Path, p)
	var input *service.PutObjectInput
	if isDir {
		cp += "/"
		input = &service.PutObjectInput{
			XQSStorageClass: convert.String(c.StorageClass),
		}
	} else {
		input = &service.PutObjectInput{
			// wrap by limitReader to keep body consistent with size
			Body:            io.LimitReader(r, size),
			ContentLength:   convert.Int64(size),
			XQSStorageClass: convert.String(c.StorageClass),
		}
	}

	var (
		contentType string
		ok          bool
	)
	if contentType, ok = meta["ContentType"]; ok {
		input.ContentType = service.String(contentType)
		delete(meta, "ContentType")
	}

	if c.UserDefineMeta {
		metadata := make(map[string]string)
		for k, v := range meta {
			metadata[strings.ToLower(k)] = v
		}

		input.XQSMetaData = &metadata
	}

	_, err = c.client.PutObject(cp, input)
	if err != nil {
		return
	}

	logrus.Debugf("QingStor wrote object %s.", cp)
	return
}

// Fetch implement destination.Fetch
func (c *Client) Fetch(ctx context.Context, p, url string) (err error) {
	cp := utils.Join(c.Path, p)

	_, err = c.client.PutObject(cp, &service.PutObjectInput{
		XQSFetchSource: convert.String(url),
	})
	if err != nil {
		return
	}

	logrus.Debugf("QingStor fetched object %s.", cp)
	return
}

// Partable implement destination.Partable
func (c *Client) Partable() bool {
	return true
}

// InitPart implement destination.InitPart
func (c *Client) InitPart(ctx context.Context, p string, size int64, meta map[string]string) (uploadID string, partSize int64, partNumbers int, err error) {
	cp := utils.Join(c.Path, p)

	input := &service.InitiateMultipartUploadInput{
		XQSStorageClass: convert.String(c.StorageClass),
	}

	var (
		contentType string
		ok          bool
	)
	if contentType, ok = meta["ContentType"]; ok {
		input.ContentType = service.String(contentType)
		delete(meta, "ContentType")
	}

	if c.UserDefineMeta {
		metadata := make(map[string]string)
		for k, v := range meta {
			metadata[strings.ToLower(k)] = v
		}

		input.XQSMetaData = &metadata
	}

	resp, err := c.client.InitiateMultipartUpload(cp, input)
	if err != nil {
		return
	}

	uploadID = *resp.UploadID
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
	cp := utils.Join(c.Path, o.Key)

	_, err = c.client.UploadMultipart(cp, &service.UploadMultipartInput{
		// wrap by limitReader to keep body consistent with size
		Body:          io.LimitReader(r, o.Size),
		ContentLength: convert.Int64(o.Size),
		UploadID:      convert.String(o.UploadID),
		PartNumber:    convert.Int(o.PartNumber),
	})
	if err != nil {
		return
	}

	logrus.Debugf("QingStor wrote partial object %s at %d.", o.Key, o.Offset)

	return nil
}

func (c *Client) CompleteParts(ctx context.Context, path string, uploadId string, totalNumber int) (err error) {
	cp := utils.Join(c.Path, path)

	logrus.Infof("Object %s start completing part", path)

	parts := make([]*service.ObjectPartType, totalNumber)
	for i := 0; i < totalNumber; i++ {
		parts[i] = &service.ObjectPartType{
			PartNumber: convert.Int(i),
		}
	}

	_, err = c.client.CompleteMultipartUpload(
		cp, &service.CompleteMultipartUploadInput{
			UploadID:    convert.String(uploadId),
			ObjectParts: parts,
		})
	if err != nil {
		return err
	}

	return nil
}
