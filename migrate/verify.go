package migrate

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/model"
)

// Verify will do verify job between src and dst.
func Verify(ctx context.Context) (err error) {
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
		go verifyWorker(ctx)
	}

	err = List(ctx)
	if err != nil {
		logrus.Errorf("List failed for %v.", err)
		return err
	}

	return
}

func verifyWorker(ctx context.Context) {
	for {
		select {
		case o, ok := <-oc:
			if !ok {
				oc = nil
				continue
			}

			err := verifyObject(ctx, o.Key)
			if err != nil {
				continue
			}
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

func verifyObject(ctx context.Context, p string) (err error) {
	defer wg.Done()

	logrus.Infof("Start verifying object %s.", p)

	so, err := src.Stat(ctx, p)
	if err != nil {
		logrus.Errorf("Src stat %s failed for %v.", p, err)
		return err
	}
	if so == nil {
		logrus.Warnf("Src object %s is not found.", p)
		return
	}

	do, err := dst.Stat(ctx, p)
	if err != nil {
		logrus.Errorf("Dst stat %s failed for %v.", p, err)
		return
	}
	// Check existence.
	if do == nil {
		logrus.Infof("Dst object %s is not found, add to repair list.", p)
		return
	}
	// Check size.
	if so.Size != do.Size {
		logrus.Infof("Dst object %s size is not match, add to repair list.", p)
		return
	}
	// Check content md5.
	sm := so.MD5
	dm := do.MD5
	if len(sm) != 32 {
		sm, err = src.MD5(ctx, p)
		if err != nil {
			logrus.Errorf("Src md5 sum failed for %v.", err)
			return
		}
		if len(sm) != 32 {
			logrus.Errorf("Src doesn't support md5 sum.")
			return
		}
	}
	if len(dm) != 32 {
		dm, err = dst.MD5(ctx, p)
		if err != nil {
			logrus.Errorf("Dst md5 sum failed for %v.", err)
			return
		}
		if len(dm) != 32 {
			logrus.Errorf("Dst doesn't support md5 sum.")
			return
		}
	}
	if sm != dm {
		logrus.Infof("Dst object %s md5 is not match, add to repair list.", p)
		return
	}

	err = model.DeleteObject(ctx, p)
	if err != nil {
		logrus.Panic(err)
	}
	return
}
