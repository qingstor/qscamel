package migrate

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

func listJob(ctx context.Context, j *model.Job) (err error) {
	defer jwg.Done()
	defer utils.Recover()

	err = src.List(ctx, j, func(o *model.Object) {
		defer utils.Recover()

		if o.IsDir {
			_, err := model.CreateJob(ctx, o.Key)
			if err != nil {
				utils.CheckClosedDB(err)
			}

			logrus.Debugf("Job %s created.", o.Key)
			return
		}

		err = model.CreateObject(ctx, o)
		if err != nil {
			utils.CheckClosedDB(err)
		}
		oc <- o
	})
	if err != nil {
		logrus.Errorf("Src list failed for %v.", err)
		return
	}

	err = model.DeleteJob(ctx, j.Path)
	if err != nil {
		utils.CheckClosedDB(err)
	}
	return
}
