package qingstor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// GetZone determines the zone of bucket.
func (q QingStor) GetZone() (zone string, err error) {
	url := fmt.Sprintf("%s://%s.%s:%s", q.Protocol, q.BucketName, q.Host, q.Port)

	r, err := http.Head(url)
	if err != nil {
		logrus.Errorf("Get QingStor zone failed for %v.", err)
		return
	}

	// Example URL: https://bucket.zone.qingstor.com
	zone = strings.Split(r.Request.URL.String(), ".")[1]
	return
}
