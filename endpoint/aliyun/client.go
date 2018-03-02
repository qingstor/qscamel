package aliyun

import (
	"context"
	"errors"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

var (
	t *model.Task
)

// Client is the client to visit aliyun oss service.
type Client struct {
	Endpoint        string
	BucketName      string
	AccessKeyID     string
	AccessKeySecret string

	Path string

	client *oss.Bucket
}

// New will create a client.
func New(ctx context.Context, et uint8) (c *Client, err error) {
	t, err = model.GetTask(ctx)
	if err != nil {
		return
	}

	c = &Client{}

	e := t.Src
	if et == constants.DestinationEndpoint {
		e = t.Dst
	}

	// Set endpoint
	c.Endpoint = e.Options["endpoint"]
	if c.Endpoint == "" {
		logrus.Error("Aliyun OSS's endpoint can't be empty.")
		err = errors.New("aliyun oss endpoint is empty")
		return
	}

	// Set bucket name.
	c.BucketName = e.Options["bucket_name"]
	if c.BucketName == "" {
		logrus.Error("Aliyun OSS's bucket name can't be empty.")
		err = errors.New("aliyun oss bucket name is empty")
		return
	}

	// Set access key.
	c.AccessKeyID = e.Options["access_key_id"]
	if c.AccessKeyID == "" {
		logrus.Error("Aliyun OSS's access key id can't be empty.")
		err = errors.New("aliyun oss access key is empty")
		return
	}

	// Set secret key.
	c.AccessKeySecret = e.Options["access_key_secret"]
	if c.AccessKeySecret == "" {
		logrus.Error("Aliyun OSS's access key secret can't be empty.")
		err = errors.New("aliyun oss access key secret is empty")
		return
	}

	// Set prefix.
	c.Path = e.Path

	service, err := oss.New(c.Endpoint, c.AccessKeyID, c.AccessKeySecret)
	if err != nil {
		return
	}
	c.client, err = service.Bucket(c.BucketName)
	if err != nil {
		return
	}

	return
}
