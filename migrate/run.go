package migrate

import (
	"context"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// Run will execute task run.
func Run(ctx context.Context) (err error) {
	switch t.Type {
	case constants.TaskTypeCopy:
		err = copyTask(ctx)
		if err != nil {
			return
		}
	case constants.TaskTypeFetch:
		err = fetchTask(ctx)
		if err != nil {
			return
		}
	case constants.TaskTypeVerifyCopy:
		err = verifyTask(ctx)
		if err != nil {
			return
		}
		err = copyTask(ctx)
		if err != nil {
			return
		}
	case constants.TaskTypeVerifyFetch:
		err = verifyTask(ctx)
		if err != nil {
			return
		}
		err = fetchTask(ctx)
		if err != nil {
			return
		}
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

func copyTask(ctx context.Context) (err error) {
	if !CanCopy() {
		logrus.Infof("Source type %s and destination type %s not support copy.",
			t.Src.Type, t.Dst.Type)
		return
	}
	logrus.Debugf("Start copy task.")

	bo := backoff.NewExponentialBackOff()

	return backoff.Retry(func() error {
		err := Copy(ctx)
		if err != nil {
			return err
		}

		if !isFinished(ctx) {
			bo.Reset()
			return constants.ErrTaskNotFinished
		}

		return nil
	}, bo)
}

func fetchTask(ctx context.Context) (err error) {
	if !CanFetch() {
		logrus.Infof("Source type %s and destination type %s not support fetch.",
			t.Src.Type, t.Dst.Type)
		return
	}
	logrus.Debugf("Start fetch task.")

	bo := backoff.NewExponentialBackOff()

	return backoff.Retry(func() error {
		err := Fetch(ctx)
		if err != nil {
			return err
		}

		if !isFinished(ctx) {
			bo.Reset()
			return constants.ErrTaskNotFinished
		}

		return nil
	}, bo)
}

func verifyTask(ctx context.Context) (err error) {
	logrus.Debugf("Start verify task.")
	return backoff.Retry(func() error {
		err = Verify(ctx)
		if err != nil {
			return err
		}

		switch t.Type {
		case constants.TaskTypeVerifyCopy:
			t.Type = constants.TaskTypeCopy
		case constants.TaskTypeVerifyFetch:
			t.Type = constants.TaskTypeFetch
		default:
			logrus.Errorf("Task %s's type %s is not supported.", t.Name, t.Type)
			return nil
		}
		err = t.Save(ctx)
		if err != nil {
			return err
		}
		return nil
	}, backoff.NewExponentialBackOff())
}

// isFinished will check whether current task has been finished.
func isFinished(ctx context.Context) bool {
	ho, err := model.HasObject(ctx)
	if err != nil {
		logrus.Panic(err)
	}
	if ho {
		logrus.Infof("There are not finished objects.")
		return false
	}

	hj, err := model.HasJob(ctx)
	if err != nil {
		logrus.Panic(err)
	}
	if hj {
		logrus.Infof("There are not finished jobs.")
		return false
	}

	return true
}
