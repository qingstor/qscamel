package qingstor

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yunify/qingstor-sdk-go/config"
	"github.com/yunify/qingstor-sdk-go/service"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

var (
	t *model.Task
)

// Client is the client to visit QingStor service.
type Client struct {
	Protocol        string
	Host            string
	Port            string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string

	Path string

	client *service.Bucket
}

// New will create a new QingStor client.
func New(ctx context.Context, et uint8) (c *Client, err error) {
	t, err = model.GetTask(ctx)
	if err != nil {
		return
	}

	c = &Client{}

	e := t.Src
	if et == constants.DestinationEndpoint {
		e = t.Dst
	}

	// Set protocol.
	c.Protocol = e.Options["protocol"]
	if c.Protocol == "" {
		c.Protocol = "https"
	}

	// Set host.
	c.Host = e.Options["host"]
	if c.Host == "" {
		c.Host = "qingstor.com"
	}

	// Set port.
	c.Port = e.Options["port"]
	if c.Port == "" {
		if c.Protocol == "https" {
			c.Port = "443"
		} else {
			c.Port = "80"
		}
	}

	// Set bucket name.
	c.BucketName = e.Options["bucket_name"]
	if c.BucketName == "" {
		logrus.Error("QingStor's bucket name can't be empty.")
		err = errors.New("qingstor bucket name is empty")
		return
	}

	// Set access key.
	c.AccessKeyID = e.Options["access_key_id"]
	if c.AccessKeyID == "" {
		logrus.Error("QingStor's access key id can't be empty.")
		err = errors.New("qingstor access key is empty")
		return
	}

	// Set secret key.
	c.SecretAccessKey = e.Options["secret_access_key"]
	if c.SecretAccessKey == "" {
		logrus.Error("QingStor's secret access key can't be empty.")
		err = errors.New("qingstor secret access key is empty")
		return
	}

	// Set path.
	c.Path = e.Path

	// Set qingstor config.
	qsConfig, _ := config.New(c.AccessKeyID, c.SecretAccessKey)
	qsConfig.Protocol = c.Protocol
	qsConfig.Host = c.Host
	port, err := strconv.ParseInt(c.Port, 10, 64)
	if err != nil {
		return
	}
	qsConfig.Port = int(port)
	// Set timeout to 0 to ignore file already closed error.
	qsConfig.Connection.Timeout = 0
	qsConfig.Connection.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			// With or without a timeout, the operating system may impose
			// its own earlier timeout
			Timeout: 1 * time.Minute,
			// Do not keep alive for too long.
			KeepAlive: 30 * time.Second,
			// XXX: DualStack enables RFC 6555-compliant "Happy Eyeballs" dialing
			// when the network is "tcp" and the destination is a host name
			// with both IPv4 and IPv6 addresses. This allows a client to
			// tolerate networks where one address family is silently broken
			DualStack: false,
		}).DialContext,
		MaxIdleConns:          0,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second, //Default
		ExpectContinueTimeout: 2 * time.Second,
	}

	qsService, _ := service.Init(qsConfig)
	zone, err := c.GetZone()
	if err != nil {
		return
	}
	c.client, _ = qsService.Bucket(c.BucketName, zone)
	return
}
