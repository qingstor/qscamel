package qingstor

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/yunify/qingstor-sdk-go/config"
	qsErrors "github.com/yunify/qingstor-sdk-go/request/errors"
	"github.com/yunify/qingstor-sdk-go/service"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Client is the client to visit QingStor service.
type Client struct {
	Protocol        string `yaml:"protocol"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	BucketName      string `yaml:"bucket_name"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`

	Path string

	client *service.Bucket
}

// New will create a new QingStor client.
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

	// Set protocol.
	if c.Protocol == "" {
		c.Protocol = "https"
	}

	// Set host.
	if c.Host == "" {
		c.Host = "qingstor.com"
	}

	// Set port.
	if c.Port == 0 {
		if c.Protocol == "https" {
			c.Port = 443
		} else {
			c.Port = 80
		}
	}

	// Set bucket name.
	if c.BucketName == "" {
		logrus.Error("QingStor's bucket name can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set access key.
	if c.AccessKeyID == "" {
		logrus.Error("QingStor's access key id can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set secret key.
	if c.SecretAccessKey == "" {
		logrus.Error("QingStor's secret access key can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set path.
	c.Path = e.Path

	// Set qingstor config.
	qc, _ := config.New(c.AccessKeyID, c.SecretAccessKey)
	qc.Protocol = c.Protocol
	qc.Host = c.Host
	qc.Port = c.Port
	qc.Connection = utils.DefaultClient

	// Set qingstor service.
	qs, _ := service.Init(qc)
	zone, err := c.GetZone()
	if err != nil {
		return
	}
	c.client, _ = qs.Bucket(c.BucketName, zone)

	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.Object, err error) {
	cp := utils.Join(c.Path, p)

	resp, err := c.client.HeadObject(cp, nil)
	if err != nil {
		if e, ok := err.(*qsErrors.QingStorError); ok {
			// If object not found, we just need to return a nil object.
			if e.StatusCode == 404 {
				return nil, nil
			}
		}
		logrus.Errorf("Stat failed for %v.", err)
		return
	}
	o = &model.Object{
		Key:          p,
		IsDir:        *resp.ContentType == DirectoryContentType,
		Size:         *resp.ContentLength,
		ContentMD5:   *resp.ETag,
		LastModified: (*resp.LastModified).Unix(),
	}
	return
}
