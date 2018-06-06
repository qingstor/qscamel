package qingstor

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/yunify/qingstor-sdk-go/client/upload"
	"github.com/yunify/qingstor-sdk-go/config"
	"github.com/yunify/qingstor-sdk-go/service"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
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

	StorageClass          string `yaml:"storage_class"`
	MultipartBoundarySize int64  `yaml:"multipart_boundary_size"`

	Path string

	client   *service.Bucket
	uploader *upload.Uploader
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
	if c.MultipartBoundarySize == 0 {
		c.MultipartBoundarySize = DefaultMultipartBoundarySize
	}
	if c.MultipartBoundarySize < 0 ||
		c.MultipartBoundarySize > MaxMultipartBoundarySize {
		logrus.Errorf("QingStor's multipart boundary size can't be %d.", c.MultipartBoundarySize)
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
	qc.Connection = hc

	// Set qingstor service.
	qs, _ := service.Init(qc)
	if c.Zone == "" {
		c.Zone, err = c.GetZone()
		if err != nil {
			return
		}
	}
	c.client, _ = qs.Bucket(c.BucketName, c.Zone)
	c.uploader = upload.Init(c.client, DefaultMultipartSize)

	return
}
