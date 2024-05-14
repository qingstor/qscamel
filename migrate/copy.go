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

// CanCopy will return whether qscamel can copy between the src and dst.
func CanCopy() bool {
	// If dst isn't writable, can't copy.
	if !dst.Writable() {
		return false
	}
	return true
}

// Copy will do copy job between src and dst.
func Copy(ctx context.Context) (err error) {
	oc = make(chan model.Object, contexts.Config.Concurrency*2)
	jc = make(chan *model.DirectoryObject)

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

// copyTask will execute a copy task.
func copyTask(ctx context.Context) (err error) {
	if !CanCopy() {
		logrus.Infof("Source type %s and destination type %s not support copy.",
			t.Src.Type, t.Dst.Type)
		return
	}
	logrus.Debugf("Start copy task.")

	bo := &backoff.ZeroBackOff{}

	return backoff.Retry(func() error {
		err := Copy(ctx)
		if err != nil {
			return err
		}

		if !isFinished(ctx) {
			//t.Status = constants.TaskStatusRerun
			return constants.ErrTaskNotFinished
		}

		return nil
	}, bo)
}
