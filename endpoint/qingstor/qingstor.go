package qingstor

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/yunify/qingstor-sdk-go/config"
	"github.com/yunify/qingstor-sdk-go/service"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"net"
	"net/http"
	"strconv"
	"time"
)

var (
	t *model.Task
)

// QingStor is the client to visit QingStor service.
type QingStor struct {
	Protocol        string
	Host            string
	Port            string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string

	Prefix string

	client *service.Bucket
}

// New will create a new QingStor client.
func New(ctx context.Context, et uint8) (q *QingStor, err error) {
	t, err = model.GetTask(ctx)
	if err != nil {
		return
	}

	q = &QingStor{}

	e := t.Src
	if et == constants.DestinationEndpoint {
		e = t.Dst
	}

	// Set protocol.
	q.Protocol = e.Options["protocol"]
	if q.Protocol == "" {
		q.Protocol = "https"
	}

	// Set host.
	q.Host = e.Options["host"]
	if q.Host == "" {
		q.Host = "qingstor.com"
	}

	// Set port.
	q.Port = e.Options["port"]
	if q.Port == "" {
		if q.Protocol == "https" {
			q.Port = "443"
		} else {
			q.Port = "80"
		}
	}

	// Set bucket name.
	q.BucketName = e.Options["bucket_name"]
	if q.BucketName == "" {
		logrus.Error("QingStor's bucket name can't be empty.")
		err = errors.New("qingstor bucket name is empty")
		return
	}

	// Set access key.
	q.AccessKeyID = e.Options["access_key_id"]
	if q.AccessKeyID == "" {
		logrus.Error("QingStor's access key id can't be empty.")
		err = errors.New("qingstor access key is empty")
		return
	}

	// Set secret key.
	q.SecretAccessKey = e.Options["secret_access_key"]
	if q.SecretAccessKey == "" {
		logrus.Error("QingStor's secret access key can't be empty.")
		err = errors.New("qingstor secret access key is empty")
		return
	}

	// Set prefix.
	q.Prefix = e.Path

	// Set qingstor config.
	qsConfig, _ := config.New(q.AccessKeyID, q.SecretAccessKey)
	qsConfig.Protocol = q.Protocol
	qsConfig.Host = q.Host
	port, err := strconv.ParseInt(q.Port, 10, 64)
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
	zone, err := q.GetZone()
	if err != nil {
		return
	}
	q.client, _ = qsService.Bucket(q.BucketName, zone)
	return
}
