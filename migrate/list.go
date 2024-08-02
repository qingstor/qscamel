package migrate

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

// List will list objects and send to channel.
func List(ctx context.Context) (err error) {
	if t.Status == constants.TaskStatusCreated {
		o := &model.DirectoryObject{
			Key:    "/",
			Marker: "",
		}
		err = model.CreateObject(ctx, o)
		if err != nil {
			logrus.Panic(err)
		}

		t.Status = constants.TaskStatusRunning
		err = t.Save(ctx)
		if err != nil {
			logrus.Panic(err)
		}

		jwg.Add(1)
		jc <- o
		return nil
	}

	// Traverse already running but not finished single object.
	p := ""
	for {
		so, err := model.NextSingleObject(ctx, p)
		if err != nil {
			logrus.Panic(err)
		}
		if so == nil {
			break
		}

		oc <- so
		p = so.Key
	}

	// Traverse already running but not finished partial object.
	p = ""
	pn := -1
	for {
		po, err := model.NextPartialObject(ctx, p, pn)
		if err != nil {
			logrus.Panic(err)
		}
		if po == nil {
			break
		}

		oc <- po
		p = po.Key
		pn = po.PartNumber
	}

	// Traverse already running but not finished directory object.
	p = ""
	for {
		do, err := model.NextDirectoryObject(ctx, p)
		if err != nil {
			logrus.Panic(err)
		}
		if do == nil {
			break
		}

		jwg.Add(1)
		jc <- do
		p = do.Key
	}

	return
}

// listWorker will do both list and copy work.
func listWorker(ctx context.Context) {
	defer utils.Recover()

	for j := range jc {
		logrus.Infof("Start listing job %s.", j.Key)

		err := listObject(ctx, j)
		if err != nil {
			logrus.Errorf("List object %s failed for %v.", j.Key, err)
			continue
		}

		logrus.Infof("Job %s listed.", j.Key)
	}
}
