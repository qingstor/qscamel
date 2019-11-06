package cos

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// Client is the client to visit aliyun oss service.
type Client struct {
	BucketURL string `yaml:"bucket_url"`
	SecretID  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`

	Path string

	client *cos.Client
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

	// Set bucket url.
	if c.BucketURL == "" {
		logrus.Error("Tencent COS's bucket url can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}
	u, err := url.Parse(c.BucketURL)
	if err != nil {
		return
	}

	// Set access key.
	if c.SecretID == "" {
		logrus.Error("Tencent COS's secret id can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set secret key.
	if c.SecretKey == "" {
		logrus.Error("Tencent COS's secret key can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set prefix.
	c.Path = e.Path
	b := &cos.BaseURL{BucketURL: u}
	c.client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  c.SecretID,
			SecretKey: c.SecretKey,
		},
		Timeout: 100 * time.Second,
	})
	return
}
