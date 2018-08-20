package qingstor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
)

// ObjectParts will store multipart upload status.
type ObjectParts struct {
	TotalParts  int `msgpack:"tps"`
	RemainParts int `msgpack:"rp"`
}

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

// calculatePartSize will calculate the object's part size.
func calculatePartSize(size int64) (partSize int64, err error) {
	partSize = DefaultMultipartSize

	for size/partSize >= int64(MaxMultipartNumber) {
		if partSize < MaxAutoMultipartSize {
			partSize = partSize << 1
			continue
		}
		// Try to adjust partSize if it is too small and account for
		// integer division truncation.
		partSize = size/int64(MaxMultipartNumber) + 1
		break
	}

	if partSize > MaxMultipartBoundarySize {
		err = constants.ErrObjectTooLarge
		return
	}

	return
}
