package fs

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// List implement source.List
// errors in List may be ignored (depend on isIgnoredErr()), to avoid infinite backoff
func (c *Client) List(ctx context.Context, j *model.DirectoryObject, fn func(o model.Object)) (err error) {
	cp := filepath.Join(c.AbsPath, j.Key)

	fi, err := os.Open(cp)
	if err != nil {
		logrus.Warnf("open <%s> failed: [%v]", cp, err)
		return err
	}
	list, err := fi.Readdir(-1)
	fi.Close()
	if err != nil {
		if isIgnoredErr(err) {
			logrus.Warnf("read dir <%s> failed: [%v], ignored", fi.Name(), err)
			return nil
		}
		return err
	}

	for _, v := range list {
		// if v is a link, and client not follow link, skip it
		if v.Mode()&os.ModeSymlink != 0 && !c.Options.EnableLinkFollow {
			continue
		}

		target, err := checkLink(v, cp)
		if err != nil {
			if isIgnoredErr(err) {
				logrus.Warnf("check link for <%s> failed: [%v], skipped", v.Name(), err)
				continue
			}
			return err
		}

		if target.IsDir() {
			o := &model.DirectoryObject{
				Key: "/" + utils.Join(j.Key, v.Name()), // always use current v's name as key
			}

			fn(o)

			continue
		}

		// Skip irregular file, such as: device, io pipe, etc.
		if !target.Mode().IsRegular() {
			logrus.Infof("target <%s> skipped because its not a regular file",
				"/"+utils.Join(cp, v.Name()))
			continue
		}

		o := &model.SingleObject{
			Key:          "/" + utils.Join(j.Key, v.Name()), // always use current v's name as key
			Size:         target.Size(),
			LastModified: target.ModTime().Unix(),
		}

		fn(o)
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

// checkLink handle a FileInfo at current path and follow link if needed
func checkLink(v os.FileInfo, cp string) (os.FileInfo, error) {
	// if v is not link, return directly
	if v.Mode()&os.ModeSymlink == 0 {
		return v, nil
	}

	// otherwise, follow the link to get the target
	tarPath, err := filepath.EvalSymlinks("/" + utils.Join(cp, v.Name()))
	if err != nil {
		return nil, err
	}
	return os.Stat(tarPath)
}

// isIgnoredErr try to check whether an error should be ignored, which will lead to not retry List()
func isIgnoredErr(err error) bool {
	if strings.Contains(err.Error(), "bad file descriptor") {
		return true
	}
	return false
}
