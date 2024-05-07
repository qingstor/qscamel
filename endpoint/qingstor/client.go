package qingstor

import (
	"context"
	"net/http"
	"time"

	"github.com/qingstor/qingstor-sdk-go/v4/config"
	"github.com/qingstor/qingstor-sdk-go/v4/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// Client is the client to visit QingStor service.
type Client struct {
	Protocol        string `yaml:"protocol"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Zone            string `yaml:"zone"`
	BucketName      string `yaml:"bucket_name"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`

	StorageClass       string `yaml:"storage_class"`
	DisableURICleaning bool   `yaml:"disable_uri_cleaning"`
	EnableVirtualStyle bool   `yaml:"enable_virtual_style"`

	// Whether to migrate custom metadata
	UserDefineMeta bool `yaml:"user_define_meta"`

	Path string

	TimeoutConfig TimeoutConfig `yaml:"timeout_config"`

	client *service.Bucket
}

type TimeoutConfig struct {
	ConnectTimeout int64 `yaml:"connect_timeout"`
	ReadTimeout    int64 `yaml:"read_timeout"`
	WriteTimeout   int64 `yaml:"write_timeout" `
}

// New will create a new QingStor client.
func New(ctx context.Context, et uint8, hc *http.Client) (c *Client, err error) {
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

	// Set storage class.
	if c.StorageClass == "" {
		c.StorageClass = StorageClassStandard
	}
	if c.StorageClass != StorageClassStandard &&
		c.StorageClass != StorageClassStandardIA {
		logrus.Errorf("QingStor's storage class can't be %s.", c.StorageClass)
		err = constants.ErrEndpointInvalid
		return
	}

	var tc = c.TimeoutConfig
	var emptyTimeoutConfig TimeoutConfig
	if tc != emptyTimeoutConfig {
		if tc.ConnectTimeout == 0 {
			tc.ConnectTimeout = utils.DEFAULT_CONN_TIMEOUT
		}
		if tc.ReadTimeout == 0 {
			tc.ReadTimeout = utils.DEFAULT_READ_TIMEOUT
		}
		if tc.WriteTimeout == 0 {
			tc.WriteTimeout = utils.DEFAULT_WRITE_TIMEOUT
		}

		connT := time.Duration(tc.ConnectTimeout) * time.Second
		readT := time.Duration(tc.ReadTimeout) * time.Second
		writeT := time.Duration(tc.WriteTimeout) * time.Second

		hc.Transport = contexts.NewTransportWithDialContext(
			contexts.Config,
			contexts.Proxy,
			utils.NewDialer(connT, readT, writeT),
		)
	}

	// Set path.
	c.Path = e.Path

	// Set qingstor config.
	qc, _ := config.New(c.AccessKeyID, c.SecretAccessKey)
	qc.Protocol = c.Protocol
	qc.Host = c.Host
	qc.Port = c.Port
	qc.Connection = hc
	qc.AdditionalUserAgent = "qscamel " + constants.Version
	qc.DisableURICleaning = c.DisableURICleaning
	qc.EnableVirtualHostStyle = c.EnableVirtualStyle

	// Set qingstor service.
	qs, _ := service.Init(qc)
	if c.Zone == "" {
		c.Zone, err = c.GetZone()
		if err != nil {
			return
		}
	}
	c.client, _ = qs.Bucket(c.BucketName, c.Zone)
	return
}
