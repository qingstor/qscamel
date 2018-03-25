package upyun

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/upyun/go-sdk/upyun"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Client is the client to visit service.
type Client struct {
	BucketName string `yaml:"bucket_name"`
	Operator   string `yaml:"operator"`
	Password   string `yaml:"password"`

	Path string

	client *upyun.UpYun
}

// New will create a new client.
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

	// Set bucket name.
	if c.BucketName == "" {
		logrus.Error("upyun bucket name can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}
	// Set operator.
	if c.Operator == "" {
		logrus.Error("upyun operator can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}
	// Set password.
	if c.Password == "" {
		logrus.Error("upyun password can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set path.
	c.Path = e.Path

	cfg := &upyun.UpYunConfig{
		Bucket:   c.BucketName,
		Operator: c.Operator,
		Password: c.Password,
	}
	c.client = upyun.NewUpYun(cfg)

	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.Object, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.GetInfo(cp)
	if err != nil {
		// If not found, upyun sdk will return an error contains "HEAD 404"
		if strings.Contains(err.Error(), "HEAD 404") {
			return nil, nil
		}
		logrus.Errorf("Get %s info failed for %v.", p, err)
		return
	}
	o = &model.Object{
		Key:          p,
		IsDir:        resp.IsDir,
		Size:         resp.Size,
		LastModified: resp.Time.Unix(),
		ContentMD5:   resp.ETag,
	}
	return
}
