package migrate

import (
	"context"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// List will list objects and send to channel.
func List(ctx context.Context) (err error) {
	if t.Status == constants.TaskStatusCreated {
		_, err = model.CreateJob(ctx, "/")
		if err != nil {
			logrus.Panic(err)
		}

		t.Status = constants.TaskStatusRunning
		err = t.Save(ctx)
		if err != nil {
			logrus.Panic(err)
		}
	}

	// Traverse already running but not finished object.
	p := ""
	for {
		o, err := model.NextObject(ctx, p)
		if err != nil {
			logrus.Panic(err)
		}
		if o == nil {
			break
		}

		oc <- o
		p = o.Key
	}

	// Traverse already running but not finished job.
	p = ""
	for {
		j, err := model.NextJob(ctx, p)
		if err != nil {
			logrus.Panic(err)
		}
		if j == nil {
			break
		}

		jwg.Add(1)
		jc <- j
		p = j.Path
	}

	return
}

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

func listJob(ctx context.Context, j *model.Job) (err error) {
	defer jwg.Done()

	err = src.List(ctx, j, func(o *model.Object) {
		if o.IsDir {
			_, err := model.CreateJob(ctx, o.Key)
			if err != nil {
				logrus.Panic(err)
			}

			logrus.Debugf("Job %s created.", o.Key)
			return
		}

		err = model.CreateObject(ctx, o)
		if err != nil {
			logrus.Panic(err)
		}
		oc <- o
	})
	if err != nil {
		logrus.Errorf("Src list failed for %v.", err)
		return
	}

	err = model.DeleteJob(ctx, j.Path)
	if err != nil {
		logrus.Panic(err)
	}
	return
}

// listWorker will do both list and copy work.
func listWorker(ctx context.Context) {
	defer utils.Recover()

	for j := range jc {
		logrus.Infof("Start list job %s.", j.Path)

		err := listJob(ctx, j)
		if err != nil {
			continue
		}

		logrus.Infof("Job %s listed.", j.Path)
	}
}

// migrateWorker will only do migrate work.
func migrateWorker(ctx context.Context) {
	defer owg.Done()
	defer utils.Recover()

	for o := range oc {
		ok, err := checkObject(ctx, o.Key)
		if err != nil || ok {
			err = model.DeleteObject(ctx, o.Key)
			if err != nil {
				logrus.Errorf("Delete object failed for %v.", err)
			}
			continue
		}

		var fn func(ctx context.Context, o *model.Object) (err error)

		switch t.Type {
		case constants.TaskTypeCopy:
			fn = copyObject
		case constants.TaskTypeFetch:
			fn = fetchObject
		default:
			logrus.Fatalf("Not supported task type: %s.", t.Type)
		}

		logrus.Infof("Start %sing object %s.", t.Type, o.Key)

		bo := backoff.NewExponentialBackOff()

		backoff.Retry(func() error {
			err = fn(ctx, o)
			if err == nil {
				return nil
			}

			logrus.Infof("Object %s %s failed, retrying.", o.Key, t.Type)
			return err
		}, bo)

		err = model.DeleteObject(ctx, o.Key)
		if err != nil {
			logrus.Errorf("Delete object failed for %v.", err)
			continue
		}

		logrus.Infof("Object %s %sed.", o.Key, t.Type)
	}
}
