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
	if t.IgnoreExisting == constants.TaskIgnoreExistingDisable {
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

	if t.IgnoreExisting == constants.TaskIgnoreExistingSize {
		logrus.Infof("Object %s check passed, ignore.", o.Key)
		return
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

	if t.IgnoreExisting != constants.TaskIgnoreExistingQuickMD5Sum &&
		t.IgnoreExisting != constants.TaskIgnoreExistingFullMD5Sum {
		return
	}

	if len(ro.MD5) != 32 {
		ro.MD5, err = md5sum(ctx, e, o)
		if err != nil {
			logrus.Errorf(
				"%s calculate object %s md5 failed for %v.", e.Name(ctx), o.Key, err)
			return
		}
		// If it's the full md5, we can update the object md5.
		if t.IgnoreExisting == constants.TaskIgnoreExistingFullMD5Sum {
			o.MD5 = ro.MD5
		}
	}
	return
}

// quickSumObject will get the object's quick md5
func quickSumObject(
	ctx context.Context, e endpoint.Base, o *model.Object,
) (m string, err error) {
	// If object size <= 3MB, use full sum instead.
	if o.Size <= 3*constants.MB {
		return fullSumObject(ctx, e, o)
	}

	goldenPoint := int64(float64(o.Size) * constants.GoldenRatio)

	pos := [][]int64{
		{0, constants.MB - 1},
		{goldenPoint, goldenPoint + constants.MB - 1},
		{o.Size - constants.MB - 1, o.Size - 1},
	}
	content := make([]byte, 3*constants.MB)

	for _, v := range pos {
		c, err := e.ReadAt(ctx, o.Key, v[0], v[1])
		if err != nil {
			logrus.Errorf("%s read object %s failed for %v.", e.Name(ctx), o.Key, err)
			return "", err
		}
		content = append(content, c...)
	}

	sum := md5.Sum(content)
	return hex.EncodeToString(sum[:]), nil
}

// fullSumObject will get the object's full md5
func fullSumObject(
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
