package migrate

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/endpoint"
	"github.com/yunify/qscamel/model"
)

// checkObject will tell whether an object is ok.
func checkObject(ctx context.Context, o *model.Object) (ok bool, err error) {
	if t.IgnoreExisting == "" {
		return false, nil
	}

	logrus.Infof("Start checking object %s.", o.Key)

	so, err := statObject(ctx, src, o)
	if err != nil {
		return
	}
	if so == nil {
		return true, nil
	}

	do, err := statObject(ctx, dst, o)
	if err != nil {
		return
	}
	// Check existence.
	if do == nil {
		return
	}
	// Check size.
	if so.Size != do.Size {
		logrus.Infof("Object %s size is not match, execute an operation on it.", o.Key)
		return
	}

	// Check last modified
	if t.IgnoreExisting == constants.TaskIgnoreExistingLastModified {
		if so.LastModified > do.LastModified {
			logrus.Infof("Object %s was modified, execute an operation on it.", o.Key)
			return
		}
		logrus.Infof("Object %s check passed, ignore.", o.Key)
		return true, nil
	}

	// Check md5.
	if so.MD5 != do.MD5 {
		logrus.Infof("Object %s md5 is not match, execute an operation on it.", o.Key)
		return
	}

	logrus.Infof("Object %s check passed, ignore.", o.Key)
	return true, nil
}

// copyObject will do a real copy.
func copyObject(ctx context.Context, o *model.Object) (err error) {
	r, err := src.Read(ctx, o.Key)
	if err != nil {
		logrus.Errorf("Src read %s failed for %v.", o.Key, err)
		return err
	}
	err = dst.Write(ctx, o.Key, o.Size, r)
	if err != nil {
		logrus.Errorf("Dst write %s failed for %v.", o.Key, err)
		return err
	}
	return
}

// deleteObject will do a real delete.
func deleteObject(ctx context.Context, o *model.Object) (err error) {
	err = dst.Delete(ctx, o.Key)
	if err != nil {
		logrus.Errorf("Dst delete %s failed for %v.", o.Key, err)
		return err
	}
	return
}

// fetchObject will do a real fetch.
func fetchObject(ctx context.Context, o *model.Object) (err error) {
	url, err := src.Reach(ctx, o.Key)
	if err != nil {
		logrus.Errorf("Src reach %s failed for %v.", o.Key, err)
		return err
	}
	err = dst.Fetch(ctx, o.Key, url)
	if err != nil {
		logrus.Errorf("Dst fetch %s failed for %v.", o.Key, err)
		return err
	}
	return
}

// statObject will get an object metadata and try to get it's md5 if available.
func statObject(
	ctx context.Context, e endpoint.Base, o *model.Object,
) (ro *model.Object, err error) {
	ro, err = e.Stat(ctx, o.Key)
	if err != nil {
		logrus.Errorf("%s stat object %s failed for %v.", e.Name(ctx), o.Key, err)
		return
	}
	if ro == nil {
		logrus.Infof("Object %s is not found at %s.", o.Key, e.Name(ctx))
		return
	}

	if t.IgnoreExisting != constants.TaskIgnoreExistingMD5Sum {
		return
	}

	if len(ro.MD5) != 32 {
		ro.MD5, err = md5SumObject(ctx, e, o)
		if err != nil {
			logrus.Errorf(
				"%s calculate object %s md5 failed for %v.", e.Name(ctx), o.Key, err)
			return
		}
	}
	return
}

// md5SumObject will get the object's md5
func md5SumObject(
	ctx context.Context, e endpoint.Base, o *model.Object,
) (m string, err error) {
	r, err := e.Read(ctx, o.Key)
	if err != nil {
		logrus.Errorf("%s read object %s failed for %v.", e.Name(ctx), o.Key, err)
		return
	}
	defer r.Close()

	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	sum := h.Sum(nil)

	return hex.EncodeToString(sum[:]), nil
}
