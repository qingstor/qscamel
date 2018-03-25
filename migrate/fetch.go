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
	oc = make(chan *model.Object)
	jc = make(chan *model.Job)
	wg = new(sync.WaitGroup)

	// Close channel for no more object.
	defer close(oc)
	// Close channel for no more job.
	defer close(jc)
	// Wait for all job finished.
	defer wg.Wait()

	for i := 0; i < contexts.Config.Concurrency; i++ {
		go fetchWorker(ctx)
	}

	err = List(ctx)
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

func fetchWorker(ctx context.Context) {
	for {
		select {
		case o, ok := <-oc:
			if !ok {
				oc = nil
				continue
			}

			if t.IgnoreExisting {
				exist, err := headObject(ctx, o.Key)
				if err != nil || exist {
					wg.Done()
					continue
				}
			}

			logrus.Infof("Start fetching object %s.", o.Key)

			err := fetchObject(ctx, o.Key)
			if err != nil {
				continue
			}

			logrus.Infof("Object %s fetched.", o.Key)
		case j, ok := <-jc:
			if !ok {
				jc = nil
				continue
			}

			logrus.Infof("Start list job %s.", j.Path)

			err := listJob(ctx, j)
			if err != nil {
				continue
			}

			logrus.Infof("Job %s listed.", j.Path)
		}
		// Check and exit while all channel closed.
		if oc == nil && jc == nil {
			break
		}
	}
}

// fetchObject will do a real fetch.
func fetchObject(ctx context.Context, p string) (err error) {
	defer wg.Done()

	url, err := src.Reach(ctx, p)
	if err != nil {
		logrus.Errorf("Src reach %s failed for %v.", p, err)
		return err
	}
	err = dst.Fetch(ctx, p, url)
	if err != nil {
		logrus.Errorf("Dst fetch %s failed for %v.", p, err)
		return err
	}

	err = model.DeleteObject(ctx, p)
	if err != nil {
		logrus.Panicf("DeleteRunningObject failed for %v.", err)
		return err
	}
	return
}
