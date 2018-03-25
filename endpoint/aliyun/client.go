package aliyun

import (
	"context"
	"strconv"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pengsrc/go-shared/convert"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Client is the client to visit aliyun oss service.
type Client struct {
	Endpoint        string `yaml:"endpoint"`
	BucketName      string `yaml:"bucket_name"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`

	Path string

	client *oss.Bucket
}

// New will create a client.
func New(ctx context.Context, et uint8) (c *Client, err error) {
	t, err := model.GetTask(ctx)
	if err != nil {
		return
	}

	c = &Client{}

	e := t.Src
	if et == constants.DestinationEndpoint {
		e = t.Dst
	}

	content, err := yaml.Marshal(e.Options)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(content, c)
	if err != nil {
		return
	}

	// Set endpoint
	if c.Endpoint == "" {
		logrus.Error("Aliyun OSS's endpoint can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set bucket name.
	if c.BucketName == "" {
		logrus.Error("Aliyun OSS's bucket name can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set access key.
	if c.AccessKeyID == "" {
		logrus.Error("Aliyun OSS's access key id can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set secret key.
	if c.AccessKeySecret == "" {
		logrus.Error("Aliyun OSS's access key secret can't be empty.")
		err = constants.ErrEndpointInvalid
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

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.Object, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.GetObjectMeta(cp)
	if err != nil {
		if e, ok := err.(*oss.ServiceError); ok {
			// If object not found, we just need to return a nil object.
			if e.StatusCode == 404 {
				return nil, nil
			}
		}
		logrus.Errorf("Stat failed for %v.", err)
		return
	}

	// Parse content length.
	size, err := strconv.ParseInt(resp.Get("Content-Length"), 10, 64)
	if err != nil {
		logrus.Errorf("Content length parsed failed for %v.", err)
		return
	}
	// Parse last modified.
	lastModified := convert.StringToUnixTimestamp(resp.Get("Last-Modified"), convert.RFC822)

	o = &model.Object{
		Key:          p,
		IsDir:        strings.HasSuffix(p, "/"),
		Size:         size,
		LastModified: lastModified,
		MD5:          resp.Get("ETag"),
	}
	return
}

// MD5 implement source.MD5
func (c *Client) MD5(ctx context.Context, p string) (b string, err error) {
	return
}
