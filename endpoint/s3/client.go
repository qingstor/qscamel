package s3

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Client is the client to visit service.
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
		err = constants.ErrEndpointInvalid
		return
	}

	// Set access key.
	if c.AccessKeyID == "" {
		logrus.Error("AWS access key id can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set secret key.
	if c.SecretAccessKey == "" {
		logrus.Error("AWS's secret access key can't be empty.")
		err = constants.ErrEndpointInvalid
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
		ContentMD5:   *resp.ETag,
		LastModified: (*resp.LastModified).Unix(),
	}
	return
}
