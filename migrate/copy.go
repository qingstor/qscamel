package migrate

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
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

	migrateWorkers := int(float64(contexts.Config.Concurrency) * constants.DefaultWorkerRatio)
	for i := 0; i < migrateWorkers; i++ {
		go migrateWorker(ctx)
	}
	for i := 0; i < contexts.Config.Concurrency-migrateWorkers; i++ {
		go listWorker(ctx)
	}

	err = List(ctx)
	if err != nil {
		logrus.Errorf("List failed for %v.", err)
		return err
	}

	return
}

// copyObject will do a real copy.
func copyObject(ctx context.Context, p string) (err error) {
	defer wg.Done()

	r, err := src.Read(ctx, p)
	if err != nil {
		logrus.Errorf("Src read %s failed for %v.", p, err)
		return err
	}
	err = dst.Write(ctx, p, r)
	if err != nil {
		logrus.Errorf("Dst write %s failed for %v.", p, err)
		return err
	}
	err = model.DeleteObject(ctx, p)
	if err != nil {
		logrus.Panicf("DeleteObject failed for %v.", err)
	}
	return
}
