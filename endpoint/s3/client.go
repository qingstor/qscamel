package s3

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// Client is the client to visit S3 service.
type Client struct {
	BucketName      string `yaml:"bucket_name"`
	Endpoint        string `yaml:"endpoint"`
	Region          string `yaml:"region"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	DisableSSL      bool   `yaml:"disable_ssl"`
	UseAccelerate   bool   `yaml:"use_accelerate"`

	Path string

	client *s3.S3
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
		logrus.Error("AWS bucket name can't be empty.")
		err = errors.New("aws bucket name is empty")
		return
	}

	// Set access key.
	if c.AccessKeyID == "" {
		logrus.Error("AWS access key id can't be empty.")
		err = errors.New("aws access key is empty")
		return
	}

	// Set secret key.
	if c.SecretAccessKey == "" {
		logrus.Error("AWS's secret access key can't be empty.")
		err = errors.New("aws secret access key is empty")
		return
	}

	// Set path.
	c.Path = e.Path

	cfg := &aws.Config{
		Credentials:     credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, ""),
		Endpoint:        aws.String(c.Endpoint),
		Region:          aws.String(c.Region),
		DisableSSL:      aws.Bool(c.DisableSSL),
		S3UseAccelerate: aws.Bool(c.UseAccelerate),
	}
	sess, err := session.NewSession(cfg)
	if err != nil {
		return
	}
	c.client = s3.New(sess)

	return
}
