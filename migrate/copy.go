package migrate

import (
	"context"
	"path"
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

	for i := 0; i < contexts.Config.Concurrency; i++ {
		wg.Add(1)
		go copyWorker(ctx, c, wg)
	}

	// Get current sequence.
	seq, err := model.GetSequence(ctx)
	if err != nil {
		logrus.Panic(err)
	}

	// Insert the first node is seq == 0.
	if seq == 0 {
		_, err = model.CreateJob(ctx, "/")
		if err != nil {
			logrus.Panic(err)
		}
		seq++
	}

	// Traverse already running but not finished job.
	err = model.ListObject(ctx, func(o *model.Object) {
		c <- o.Key
	})
	if err != nil {
		logrus.Panic(err)
	}

	// Get current job IDs.
	cur, err := model.GetCurrentJobID(ctx)
	if err != nil {
		logrus.Panic(err)
	}

	for {
		if seq == cur {
			break
		}

		j, err := model.GetJob(ctx, cur+1)
		if err != nil {
			logrus.Panic(err)
		}
		// Create folder before walking.
		err = dst.Dir(ctx, j.Path)
		if err != nil {
			return err
		}

		fi, err := src.List(ctx, j.Path)
		if err != nil {
			return err
		}
		for k, v := range fi {
			if v.IsDir {
				_, err = model.CreateJob(ctx, v.Key)
				if err != nil {
					return err
				}
				// Update bucket sequence.
				seq++

				logrus.Debugf("Job %s is created.", path.Join(j.Path, v.Key))
				continue
			}

			err = model.CreateObject(ctx, &fi[k])
			if err != nil {
				return err
			}

			c <- v.Key
		}

		// Update current running job.
		cur++
		err = model.UpdateCurrentJobID(ctx, cur)
		if err != nil {
			logrus.Panic(err)
		}

		logrus.Debugf("Job %s is finished.", j.Path)
	}

	// Close channel for no more job.
	close(c)

	// Wait for all job finished.
	wg.Wait()

	ho, err := model.HasObject(ctx)
	if err != nil {
		logrus.Panic(err)
	}
	if ho {
		logrus.Infof("There are not finished objects, retried.")
		return Copy(ctx)
	}
	// TODO: we should retry failed task.
	logrus.Infof("Task %s has been finished.", t.Name)
	return
}

func copyWorker(ctx context.Context, c chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range c {
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
