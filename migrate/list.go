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
		_, err = model.CreateJob(ctx, "/")
		if err != nil {
			logrus.Panic(err)
		}

		t.Status = constants.TaskStatusRunning
		err = t.Save(ctx)
		if err != nil {
			logrus.Panic(err)
		}
	}

	// Traverse already running but not finished object.
	p := ""
	for {
		o, err := model.NextObject(ctx, p)
		if err != nil {
			logrus.Panic(err)
		}
		if o == nil {
			break
		}

		oc <- o
		p = o.Key
	}

	// Traverse already running but not finished job.
	p = ""
	for {
		j, err := model.NextJob(ctx, p)
		if err != nil {
			logrus.Panic(err)
		}
		if j == nil {
			break
		}

		jwg.Add(1)
		jc <- j
		p = j.Path
	}

	return
}

// listWorker will do both list and copy work.
func listWorker(ctx context.Context) {
	defer utils.Recover()

	for j := range jc {
		logrus.Infof("Start listing job %s.", j.Path)

		err := listJob(ctx, j)
		if err != nil {
			continue
		}

		logrus.Infof("Job %s listed.", j.Path)
	}
}
