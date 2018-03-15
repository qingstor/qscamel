package migrate

import (
	"context"
	"errors"
	"sync"

	"github.com/sirupsen/logrus"

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
	c := make(chan string)
	wg := new(sync.WaitGroup)

	// Close channel for no more job.
	defer close(c)
	// Wait for all job finished.
	defer wg.Wait()

	for i := 0; i < contexts.Config.Concurrency; i++ {
		wg.Add(1)
		go fetchWorker(ctx, c, wg)
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

func fetchWorker(ctx context.Context, c chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range c {
		logrus.Infof("Start fetching object %s.", p)

		url, err := src.Reach(ctx, p)
		if err != nil {
			logrus.Errorf("Src reach %s failed for %v.", p, err)
			continue
		}
		err = dst.Fetch(ctx, p, url)
		if err != nil {
			logrus.Errorf("Dst fetch %s failed for %v.", p, err)
			continue
		}
		err = model.DeleteObject(ctx, p)
		if err != nil {
			logrus.Panicf("DeleteRunningObject failed for %v.", err)
			continue
		}

		logrus.Infof("Object %s fetched.", p)
	}
}
