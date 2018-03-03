package gcs

import (
	"context"
	"errors"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// Client is the client to visit service.
type Client struct {
	APIKey     string `yaml:"api_key"`
	BucketName string `yaml:"bucket_name"`

	Path string

	client *storage.BucketHandle
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
		logrus.Error("Google cloud storage bucket name can't be empty.")
		err = errors.New("google cloud storage bucket name is empty")
		return
	}
	// Set api key
	if c.APIKey == "" {
		logrus.Error("Google cloud storage API key can't be empty.")
		err = errors.New("google cloud storage api key is empty")
		return
	}

	// Set path.
	c.Path = e.Path

	svc, err := storage.NewClient(ctx, option.WithAPIKey(c.APIKey))
	if err != nil {
		return
	}
	c.client = svc.Bucket(c.BucketName)
	return
}
