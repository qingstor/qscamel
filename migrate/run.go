package migrate

import (
	"context"
	"errors"
	"path"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// Run will execute task run.
func Run(ctx context.Context) (err error) {
	switch t.Type {
	case constants.TaskTypeCopy:
		if !CanCopy() {
			logrus.Infof("Source type %s and destination type %s not support copy.",
				t.Src.Type, t.Dst.Type)
			return
		}
		err = backoff.Retry(func() error { return Copy(ctx) }, backoff.NewExponentialBackOff())
		if err != nil {
			return
		}
	case constants.TaskTypeFetch:
		if !CanFetch() {
			logrus.Infof("Source type %s and destination type %s not support fetch.",
				t.Src.Type, t.Dst.Type)
			return
		}
		err = backoff.Retry(func() error { return Fetch(ctx) }, backoff.NewExponentialBackOff())
		if err != nil {
			return
		}
	case constants.TaskTypeVerify:
		return
	default:
		logrus.Errorf("Task %s's type %s is not supported.", t.Name, t.Type)
		return
	}

	// Update task status.
	t.Status = constants.TaskStatusFinished
	err = t.Save(ctx)
	if err != nil {
		logrus.Print(err)
	}

	return
}

// List will list objects and send to channel.
func List(ctx context.Context, c chan string) (err error) {
	// Get current sequence.
	seq, err := model.GetSequence(ctx)
	if err != nil {
		logrus.Panic(err)
	}

	// Insert the first node is seq == 0.
	if seq == 0 {
		_, err = model.CreateJob(ctx, "/")
		if err != nil {
			logrus.Panic(err)
		}
		seq++
	}

	// Traverse already running but not finished job.
	err = model.ListObject(ctx, func(o *model.Object) {
		c <- o.Key
	})
	if err != nil {
		logrus.Panic(err)
	}

	// Get current job IDs.
	cur, err := model.GetCurrentJobID(ctx)
	if err != nil {
		logrus.Panic(err)
	}
	for {
		if seq <= cur {
			break
		}

		j, err := model.GetJob(ctx, cur+1)
		if err != nil {
			logrus.Panic(err)
		}
		// Create folder before walking.
		err = dst.Dir(ctx, j.Path)
		if err != nil {
			return err
		}

		rc := make(chan *model.Object)

		go src.List(ctx, j, rc)

		for v := range rc {
			if v == nil {
				logrus.Errorf("Something error happened while listing.")
				err = errors.New("source list failed")
				return err
			}

			if v.IsDir {
				_, err = model.CreateJob(ctx, v.Key)
				if err != nil {
					return err
				}
				// Update bucket sequence.
				seq++

				logrus.Debugf("Job %s is created.", path.Join(j.Path, v.Key))
				continue
			}

			err = model.CreateObject(ctx, v)
			if err != nil {
				return err
			}

			c <- v.Key
		}

		// Update current running job.
		cur++
		err = model.UpdateCurrentJobID(ctx, cur)
		if err != nil {
			logrus.Panic(err)
		}

		logrus.Debugf("Job %s is finished.", j.Path)
	}

	return
}
