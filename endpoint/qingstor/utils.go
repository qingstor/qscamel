package qingstor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// GetZone determines the zone of bucket.
func (c *Client) GetZone() (zone string, err error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url := fmt.Sprintf("%s://%s.%s:%d", c.Protocol, c.BucketName, c.Host, c.Port)

	r, err := client.Head(url)
	if err != nil {
		logrus.Errorf("Get QingStor zone failed for %v.", err)
		return
	}

	// Example URL: https://bucket.zone.qingstor.com
	zone = strings.Split(r.Header.Get("Location"), ".")[1]
	return
}
