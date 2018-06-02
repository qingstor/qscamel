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

// CanFetch will return whether qscamel can fetch between the src and dst.
func CanFetch() bool {
	// If src isn't reachable, can't fetch.
	if !src.Reachable() {
		return false
	}
	// If dst isn't fetchable, can't fetch.
	if !dst.Fetchable() {
		return false
	}
	return true
}

// Fetch will do fetch job between src and dst.
func Fetch(ctx context.Context) (err error) {
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

// fetchTask will excuate a fetch task.
func fetchTask(ctx context.Context) (err error) {
	if !CanFetch() {
		logrus.Infof("Source type %s and destination type %s not support fetch.",
			t.Src.Type, t.Dst.Type)
		return
	}
	logrus.Debugf("Start fetch task.")

	bo := &backoff.ZeroBackOff{}

	return backoff.Retry(func() error {
		err := Fetch(ctx)
		if err != nil {
			return err
		}

		if !isFinished(ctx) {
			return constants.ErrTaskNotFinished
		}

		return nil
	}, bo)
}
