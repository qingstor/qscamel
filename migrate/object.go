package migrate

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
)

// checkObject will tell whether an object is ok.
func checkObject(ctx context.Context, p string) (ok bool, err error) {
	if !t.IgnoreExisting {
		return false, nil
	}

	logrus.Infof("Start checking object %s.", p)

	so, err := src.Stat(ctx, p)
	if err != nil {
		logrus.Errorf("Src stat %s failed for %v.", p, err)
		return
	}
	if so == nil {
		logrus.Warnf("Src object %s is not found, ignored.", p)
		return true, nil
	}

	do, err := dst.Stat(ctx, p)
	if err != nil {
		logrus.Errorf("Dst stat %s failed for %v.", p, err)
		return
	}
	// Check existence.
	if do == nil {
		logrus.Infof("Dst object %s is not found, should execute an operation on it.", p)
		return
	}
	// Check size.
	if so.Size != do.Size {
		logrus.Infof("Dst object %s size is not match, should execute an operation on it.", p)
		return
	}
	// Check content md5.
	if src.MD5able() && dst.MD5able() {
		sm := so.MD5
		dm := do.MD5
		if len(sm) != 32 {
			sm, err = src.MD5(ctx, p)
			if err != nil {
				logrus.Errorf("Src md5 sum failed for %v.", err)
				return
			}
		}
		if len(dm) != 32 {
			dm, err = dst.MD5(ctx, p)
			if err != nil {
				logrus.Errorf("Dst md5 sum failed for %v.", err)
				return
			}
		}
		if sm != dm {
			logrus.Infof("Dst object %s md5 is not match, should execute an operation on it.", p)
			return
		}
	}

	logrus.Infof("Object %s check passed, ignore.", p)
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
