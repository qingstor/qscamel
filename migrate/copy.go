package migrate

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/model"
)

// CanCopy will return whether qscamel can copy between the src and dst.
func CanCopy() bool {
	// If src neither reachable nor readable, can't copy.
	if !src.Reachable() && !src.Readable() {
		return false
	}
	// If dst isn't writable, can't copy.
	if !dst.Writable() {
		return false
	}
	return true
}

// Copy will do copy job between src and dst.
func Copy(ctx context.Context) (err error) {
	oc = make(chan *model.Object)
	jc = make(chan *model.Job)
	wg = new(sync.WaitGroup)

	// Close channel for no more object.
	defer close(oc)
	// Close channel for no more job.
	defer close(jc)
	// Wait for all job finished.
	defer wg.Wait()

	go listWorker(ctx)

	for i := 0; i < contexts.Config.Concurrency; i++ {
		go migrateWorker(ctx)
	}

	err = List(ctx)
	if err != nil {
		logrus.Errorf("List failed for %v.", err)
		return err
	}

	return
}

// copyObject will do a real copy.
func copyObject(ctx context.Context, o *model.Object) (err error) {
	defer wg.Done()
	defer func() {
		if err != nil {
			model.CreateObject(ctx, o)
		}
	}()

	r, err := src.Read(ctx, o.Key)
	if err != nil {
		logrus.Errorf("Src read %s failed for %v.", o.Key, err)
		return err
	}
	err = dst.Write(ctx, o.Key, r)
	if err != nil {
		logrus.Errorf("Dst write %s failed for %v.", o.Key, err)
		return err
	}
	return
}
