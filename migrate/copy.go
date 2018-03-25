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
	oc := make(chan *model.Object)
	jc := make(chan *model.Job)
	wg := new(sync.WaitGroup)

	// Close channel for no more object.
	defer close(oc)
	// Close channel for no more job.
	defer close(jc)
	// Wait for all job finished.
	defer wg.Wait()

	for i := 0; i < contexts.Config.Concurrency; i++ {
		go copyWorker(ctx, oc, jc, wg)
	}

	err = List(ctx, oc, jc, wg)
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
		return constants.ErrNotFinishedObject
	}

	return
}

func copyWorker(ctx context.Context, oc chan *model.Object, jc chan *model.Job, wg *sync.WaitGroup) {
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

			logrus.Infof("Start copying object %s.", o.Key)

			err := copyObject(ctx, o.Key, wg)
			if err != nil {
				continue
			}

			logrus.Infof("Object %s copied.", o.Key)
		case j, ok := <-jc:
			if !ok {
				jc = nil
				continue
			}

			logrus.Infof("Start list job %s.", j.Path)

			err := listJob(ctx, j, oc, jc, wg)
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

// copyObject will do a real copy.
func copyObject(ctx context.Context, p string, wg *sync.WaitGroup) (err error) {
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
