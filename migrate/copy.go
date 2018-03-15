package migrate

import (
	"context"
	"errors"
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
	c := make(chan string)
	wg := new(sync.WaitGroup)

	// Close channel for no more job.
	defer close(c)
	// Wait for all job finished.
	defer wg.Wait()

	for i := 0; i < contexts.Config.Concurrency; i++ {
		wg.Add(1)
		go copyWorker(ctx, c, wg)
	}

	err = List(ctx, c)
	if err != nil {
		logrus.Errorf("List failed for %v.", err)
		return err
	}

	ho, err := model.HasObject(ctx)
	if err != nil {
		logrus.Panic(err)
	}
	if ho {
		logrus.Infof("There are not finished objects, retried.")
		err = errors.New("object not finished")
		return
	}

	logrus.Infof("Task %s has been finished.", t.Name)
	return
}

func copyWorker(ctx context.Context, c chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range c {
		logrus.Infof("Start copying object %s.", p)

		r, err := src.Read(ctx, p)
		if err != nil {
			logrus.Errorf("Src read %s failed for %v.", p, err)
			continue
		}
		err = dst.Write(ctx, p, r)
		if err != nil {
			logrus.Errorf("Dst write %s failed for %v.", p, err)
			continue
		}
		err = model.DeleteObject(ctx, p)
		if err != nil {
			logrus.Panicf("DeleteRunningObject failed for %v.", err)
			continue
		}

		logrus.Infof("Object %s copied.", p)
	}
}
