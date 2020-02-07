package hdfs

import (
	"context"
	"net/http"

	"github.com/colinmarc/hdfs/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// Client is the struct for HDFS endpoint.
type Client struct {
	Address string `yaml:"address"`

	Path string

	client *hdfs.Client
}

// New will create a client.
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

	// Check address.
	if c.Address == "" {
		logrus.Error("HDFS's address can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	c.Path = e.Path
	c.client, err = hdfs.New(c.Address)
	if err != nil {
		return nil, err
	}
	return
}
