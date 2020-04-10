package azblob

import (
	"context"
	"net/http"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/services/azblob"
	"github.com/Xuanwo/storage/types/pairs"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// Client is the client to visit service.
type Client struct {
	AccountName string `yaml:"account_name"`
	AccountKey  string `yaml:"account_key"`
	BucketName  string `yaml:"bucket_name"`

	Path string

	client storage.Storager
}

// New will create a new client.
func New(ctx context.Context, et uint8, _ *http.Client) (c *Client, err error) {
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
		logrus.Error("Azure blob storage bucket name can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}
	// Set account
	if c.AccountName == "" {
		logrus.Error("Azure blob storage account name can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}
	if c.AccountKey == "" {
		logrus.Error("Azure blob storage account key can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set path.
	c.Path = e.Path

	_, c.client, err = azblob.New(
		pairs.WithCredential(credential.MustNewHmac(c.AccountName, c.AccountKey)),
		pairs.WithName(c.BucketName),
		pairs.WithWorkDir(c.Path),
	)
	return
}
