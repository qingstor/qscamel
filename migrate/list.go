package migrate

import (
	"context"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
)

// List will list objects and send to channel.
func List(ctx context.Context) (err error) {
	seq, err := model.GetSequence(ctx)
	if err != nil {
		logrus.Panic(err)
	}
	if seq == 0 {
		_, err = model.CreateJob(ctx, "/")
		if err != nil {
			logrus.Panic(err)
		}
	}

	// Traverse already running but not finished object.
	err = model.ListObject(ctx, func(o *model.Object) {
		wg.Add(1)
		oc <- o
	})
	if err != nil {
		logrus.Panic(err)
	}

	// Traverse already running but not finished job.
	err = model.ListJob(ctx, func(j *model.Job) {
		wg.Add(1)
		jc <- j
	})
	if err != nil {
		logrus.Panic(err)
	}

	return
}

// headObject will head an object.
func headObject(ctx context.Context, p string) (exist bool, err error) {
	exist = false

	so, err := src.Stat(ctx, p)
	if err != nil {
		logrus.Errorf("Src stat %s failed for %v.", p, err)
		return
	}
	if so == nil {
		logrus.Warnf("Src object %s is not found.", p)
		return
	}

	do, err := dst.Stat(ctx, p)
	if err != nil {
		logrus.Errorf("Dst stat %s failed for %v.", p, err)
		return
	}
	if do == nil {
		logrus.Warnf("Dst object %s is not found, execute.", p)
		return
	}

	exist = true
	logrus.Warnf("Dst object %s exists, ignore.", p)

	err = model.DeleteObject(ctx, p)
	if err != nil {
		logrus.Panicf("DeleteRunningObject failed for %v.", err)
	}
	return
}

func listJob(ctx context.Context, j *model.Job) (err error) {
	defer wg.Done()

	err = src.List(ctx, j, func(o *model.Object) {
		if o.IsDir {
			nj, err := model.CreateJob(ctx, o.Key)
			if err != nil {
				// Panic a db error
				logrus.Panic(err)
			}

			logrus.Infof("Job %s is created.", o.Key)

			wg.Add(1)
			select {
			case jc <- nj:
			default:
				wg.Done()
				break
			}

			return
		}

		err = model.CreateObject(ctx, o)
		if err != nil {
			logrus.Panic(err)
		}

		wg.Add(1)
		select {
		case oc <- o:
		default:
			wg.Done()
			break
		}
	})
	if err != nil {
		logrus.Errorf("Src list failed for %v.", err)
		return
	}

	err = model.DeleteJob(ctx, j.ID)
	if err != nil {
		logrus.Panic(err)
	}
	return
}
