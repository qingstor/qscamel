package migrate

import (
	"context"
	"sync"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/model"
)

// CanDelete will return whether qscamel can delete between the src and dst.
func CanDelete() bool {
	// If dst isn't writable, can't copy.
	if !dst.Writable() {
		return false
	}
	return true
}

// Delete will do delete job between src and dst.
func Delete(ctx context.Context) (err error) {
	oc = make(chan *model.Object, contexts.Config.Concurrency*2)
	jc = make(chan *model.Job)

	owg = &sync.WaitGroup{}
	jwg = &sync.WaitGroup{}

	// Wait for all object finished.
	defer owg.Wait()
	// Close channel for no more object.
	defer close(oc)
	// Close channel for no more job.
	defer close(jc)
	// Wait for all job finished.
	defer jwg.Wait()

	go listWorker(ctx)

	for i := 0; i < contexts.Config.Concurrency; i++ {
		owg.Add(1)
		go migrateWorker(ctx)
	}

	err = List(ctx)
	if err != nil {
		logrus.Errorf("List failed for %v.", err)
		return err
	}

	return
}

// deleteTask will execute a delete task.
func deleteTask(ctx context.Context) (err error) {
	if !CanCopy() {
		logrus.Infof("Source type %s and destination type %s not support delete.",
			t.Src.Type, t.Dst.Type)
		return
	}
	logrus.Debugf("Start delete task.")

	bo := &backoff.ZeroBackOff{}

	return backoff.Retry(func() error {
		err := Delete(ctx)
		if err != nil {
			return err
		}

		if !isFinished(ctx) {
			return constants.ErrTaskNotFinished
		}

		return nil
	}, bo)
}
