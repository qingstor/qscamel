package migrate

import (
	"context"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
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

	logrus.Infof("Task %s has been finished.", t.Name)
	return
}
