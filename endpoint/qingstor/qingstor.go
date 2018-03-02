package qingstor

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/yunify/qingstor-sdk-go/config"
	"github.com/yunify/qingstor-sdk-go/service"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
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

	logrus.Debug(e.Options)
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

	qsConfig, _ := config.New(q.AccessKeyID, q.SecretAccessKey)
	qsService, _ := service.Init(qsConfig)
	zone, err := q.GetZone()
	if err != nil {
		return
	}
	q.client, _ = qsService.Bucket(q.BucketName, zone)
	return
}
