package filelist

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// List implement source.List
func (c *Client) List(ctx context.Context, j *model.DirectoryObject, fn func(o model.Object)) (err error) {
	lp, err := filepath.Abs(c.ListPath)
	if err != nil {
		return
	}

	marker := j.Marker
	if len(marker) == 0 {
		marker = "0"
	}
	cur, err := strconv.ParseInt(marker, 10, 64)
	if err != nil {
		return
	}

	fi, err := os.Open(lp)
	if err != nil {
		return
	}
	defer fi.Close()

	buf := bufio.NewScanner(fi)

	for buf.Scan() {
		line := buf.Text()

		cur += int64(len(buf.Bytes()))

		// TODO: we should get file's size here.
		o := &model.SingleObject{
			Key: "/" + utils.Join(j.Key, line),
		}

		fn(o)

		j.Marker = strconv.FormatInt(cur, 10)
		err = model.CreateObject(ctx, j)
		if err != nil {
			logrus.Errorf("Save directory object %s failed for %v.", j.Key, err)
			return err
		}
	}
	return
}

// Reach implement source.Fetch
func (c *Client) Reach(ctx context.Context, p string) (url string, err error) {
	return "", constants.ErrEndpointFuncNotImplemented
}

// Reachable implement source.Reachable
func (c *Client) Reachable() bool {
	return false
}
