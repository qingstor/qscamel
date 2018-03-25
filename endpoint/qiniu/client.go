package qiniu

import (
	"context"
	"net/http"
	"strings"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/x/rpc.v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Client is the client to visit aliyun oss service.
type Client struct {
	BucketName    string `yaml:"bucket_name"`
	AccessKey     string `yaml:"access_key"`
	SecretKey     string `yaml:"secret_key"`
	Domain        string `yaml:"domain"`
	UseHTTPS      bool   `yaml:"use_https"`
	UseCdnDomains bool   `yaml:"use_cdn_domains"`

	Path string

	client *http.Client
	bucket *storage.BucketManager
	mac    *qbox.Mac
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
		logrus.Error("Qiniu's bucket name can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set access key.
	if c.AccessKey == "" {
		logrus.Error("Qiniu's access key can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set secret key.
	if c.SecretKey == "" {
		logrus.Error("Qiniu's secret key can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set domain.
	if c.Domain == "" {
		logrus.Error("Qiniu's domain can't be empty.")
		err = constants.ErrEndpointInvalid
		return
	}

	// Set prefix.
	c.Path = e.Path

	// Set qiniu related clients.
	c.mac = qbox.NewMac(c.AccessKey, c.SecretKey)
	zone, err := storage.GetZone(c.AccessKey, c.BucketName)
	if err != nil {
		return
	}
	cfg := &storage.Config{
		Zone:          zone,
		UseHTTPS:      c.UseHTTPS,
		UseCdnDomains: c.UseCdnDomains,
	}
	c.bucket = storage.NewBucketManager(c.mac, cfg)
	c.client = utils.DefaultClient

	return
}

// Stat implement source.Stat and destination.Stat
func (c *Client) Stat(ctx context.Context, p string) (o *model.Object, err error) {
	cp := utils.Join(c.Path, p)

	fi, err := c.bucket.Stat(c.BucketName, cp)
	if err != nil {
		if e, ok := err.(*rpc.ErrorInfo); ok {
			// If object not found, we just need to return a nil object.
			if e.Code == ErrorCodeNotFound {
				return nil, nil
			}
		}
		logrus.Errorf("Stat failed for %v.", err)
		return
	}
	// qiniu use their own hash algorithm instead of md5, so we can't support it.
	o = &model.Object{
		Key:          p,
		IsDir:        strings.HasSuffix(p, "/"),
		Size:         fi.Fsize,
		LastModified: fi.PutTime,
	}
	return
}

// MD5 implement source.MD5 and destination.MD5
func (c *Client) MD5(ctx context.Context, p string) (b string, err error) {
	return
}
