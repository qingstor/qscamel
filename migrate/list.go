package migrate

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// List will list objects and send to channel.
func List(ctx context.Context) (err error) {
	seq, err := model.GetSequence(ctx)
	if err != nil {
		logrus.Panic(err)
	}
	if seq == 0 {
		_, err = model.CreateJob(ctx, "/")
		if err != nil {
			logrus.Panic(err)
		}
	}

	// Traverse already running but not finished object.
	err = model.ListObject(ctx, func(o *model.Object) {
		wg.Add(1)
		oc <- o
	})
	if err != nil {
		logrus.Panic(err)
	}

	// Traverse already running but not finished job.
	err = model.ListJob(ctx, func(j *model.Job) {
		wg.Add(1)
		jc <- j
	})
	if err != nil {
		logrus.Panic(err)
	}

	return
}

// checkObject will tell whether an object is ok.
func checkObject(ctx context.Context, p string) (ok bool, err error) {
	if !t.IgnoreExisting {
		return true, nil
	}

	logrus.Infof("Start checking object %s.", p)

	so, err := src.Stat(ctx, p)
	if err != nil {
		logrus.Errorf("Src stat %s failed for %v.", p, err)
		return
	}
	if so == nil {
		logrus.Warnf("Src object %s is not found, ignored.", p)
		return true, nil
	}

	do, err := dst.Stat(ctx, p)
	if err != nil {
		logrus.Errorf("Dst stat %s failed for %v.", p, err)
		return
	}
	// Check existence.
	if do == nil {
		logrus.Infof("Dst object %s is not found, should execute an operation on it.", p)
		return
	}
	// Check size.
	if so.Size != do.Size {
		logrus.Infof("Dst object %s size is not match, should execute an operation on it.", p)
		return
	}
	// Check content md5.
	if src.MD5able() && dst.MD5able() {
		sm := so.MD5
		dm := do.MD5
		if len(sm) != 32 {
			sm, err = src.MD5(ctx, p)
			if err != nil {
				logrus.Errorf("Src md5 sum failed for %v.", err)
				return
			}
		}
		if len(dm) != 32 {
			dm, err = dst.MD5(ctx, p)
			if err != nil {
				logrus.Errorf("Dst md5 sum failed for %v.", err)
				return
			}
		}
		if sm != dm {
			logrus.Infof("Dst object %s md5 is not match, should execute an operation on it.", p)
			return
		}
	}

	// Delete this object if every thing check passed.
	err = model.DeleteObject(ctx, p)
	if err != nil {
		logrus.Panicf("DeleteObject failed for %v.", err)
	}
	return true, nil
}

func listJob(ctx context.Context, j *model.Job) (err error) {
	defer wg.Done()

	err = src.List(ctx, j, func(o *model.Object) {
		if o.IsDir {
			_, err := model.CreateJob(ctx, o.Key)
			if err != nil {
				// Panic a db error
				logrus.Panic(err)
			}

			logrus.Infof("Job %s is created.", o.Key)
			return
		}

		err = model.CreateObject(ctx, o)
		if err != nil {
			logrus.Panic(err)
		}

		wg.Add(1)
		oc <- o
		logrus.Infof("Object %s is created.", o.Key)
	})
	if err != nil {
		logrus.Errorf("Src list failed for %v.", err)
		return
	}

	err = model.DeleteJob(ctx, j.ID)
	if err != nil {
		logrus.Panic(err)
	}
	return
}

// listWorker will do both list and copy work.
func listWorker(ctx context.Context) {
	for {
		select {
		case o, ok := <-oc:
			if !ok {
				oc = nil
				continue
			}

			ok, err := checkObject(ctx, o.Key)
			if err != nil || ok {
				wg.Done()
				continue
			}

			switch t.Type {
			case constants.TaskTypeCopy:
				logrus.Infof("Start copying object %s.", o.Key)
				err = copyObject(ctx, o.Key)
				if err != nil {
					continue
				}

				logrus.Infof("Object %s copied.", o.Key)
			case constants.TaskTypeFetch:
				logrus.Infof("Start fetching object %s.", o.Key)

				err = fetchObject(ctx, o.Key)
				if err != nil {
					continue
				}

				logrus.Infof("Object %s fetched.", o.Key)
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

// migrateWorker will only do migrate work.
func migrateWorker(ctx context.Context) {
	for o := range oc {
		ok, err := checkObject(ctx, o.Key)
		if err != nil || ok {
			wg.Done()
			continue
		}

		switch t.Type {
		case constants.TaskTypeCopy:
			logrus.Infof("Start copying object %s.", o.Key)
			err = copyObject(ctx, o.Key)
			if err != nil {
				continue
			}

			logrus.Infof("Object %s copied.", o.Key)
		case constants.TaskTypeFetch:
			logrus.Infof("Start fetching object %s.", o.Key)

			err = fetchObject(ctx, o.Key)
			if err != nil {
				continue
			}

			logrus.Infof("Object %s fetched.", o.Key)
		}
	}
}
