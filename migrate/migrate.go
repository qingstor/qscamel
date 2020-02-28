// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package migrate

import (
	"context"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/endpoint"
	"github.com/yunify/qscamel/endpoint/aliyun"
	"github.com/yunify/qscamel/endpoint/cos"
	"github.com/yunify/qscamel/endpoint/filelist"
	"github.com/yunify/qscamel/endpoint/fs"
	"github.com/yunify/qscamel/endpoint/gcs"
	"github.com/yunify/qscamel/endpoint/hdfs"
	"github.com/yunify/qscamel/endpoint/qingstor"
	"github.com/yunify/qscamel/endpoint/qiniu"
	"github.com/yunify/qscamel/endpoint/s3"
	"github.com/yunify/qscamel/endpoint/upyun"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

var (
	t *model.Task

	oc chan model.Object
	jc chan *model.DirectoryObject

	owg *sync.WaitGroup
	jwg *sync.WaitGroup

	src endpoint.Source
	dst endpoint.Destination

	multipartBoundarySize int64
)

// Execute will execute migrate task.
func Execute(ctx context.Context) (err error) {
	t, err = model.GetTask(ctx)
	if err != nil {
		return
	}

	// If multipart boundary size is 0 or invalid, qscamel will correct it
	// to default boundary size.
	if t.MultipartBoundarySize > 0 {
		multipartBoundarySize = t.MultipartBoundarySize
	} else {
		multipartBoundarySize = constants.DefaultMultipartBoundarySize
	}

	err = check(ctx)
	if err != nil {
		logrus.Errorf("Pre migrate check failed for %v.", err)
		return
	}

	return run(ctx)
}

func check(ctx context.Context) (err error) {
	// Initialize source.
	switch t.Src.Type {
	case constants.EndpointAliyun:
		src, err = aliyun.New(ctx, constants.SourceEndpoint, contexts.Client)
		if err != nil {
			return
		}
	case constants.EndpointFileList:
		src, err = filelist.New(ctx, constants.SourceEndpoint)
		if err != nil {
			return
		}
	case constants.EndpointFs:
		src, err = fs.New(ctx, constants.SourceEndpoint)
		if err != nil {
			return
		}
	case constants.EndpointGCS:
		src, err = gcs.New(ctx, constants.SourceEndpoint, contexts.Client)
		if err != nil {
			return
		}
	case constants.EndpointHDFS:
		src, err = hdfs.New(ctx, constants.SourceEndpoint, contexts.Client)
		if err != nil {
			return
		}
	case constants.EndpointQingStor:
		src, err = qingstor.New(ctx, constants.SourceEndpoint, contexts.Client)
		if err != nil {
			return
		}
	case constants.EndpointQiniu:
		src, err = qiniu.New(ctx, constants.SourceEndpoint, contexts.Client)
		if err != nil {
			return
		}
	case constants.EndpointS3:
		src, err = s3.New(ctx, constants.SourceEndpoint, contexts.Client)
		if err != nil {
			return
		}
	case constants.EndpointUpyun:
		src, err = upyun.New(ctx, constants.SourceEndpoint, contexts.Client)
		if err != nil {
			return
		}
	case constants.EndpointCOS:
		src, err = cos.New(ctx, constants.SourceEndpoint, contexts.Client)
		if err != nil {
			return
		}
	default:
		logrus.Errorf("Type %s is not supported.", t.Src.Type)
		err = constants.ErrEndpointNotSupported
		return
	}

	// Initialize destination.
	switch t.Dst.Type {
	case constants.EndpointQingStor:
		dst, err = qingstor.New(ctx, constants.DestinationEndpoint, contexts.Client)
		if err != nil {
			return
		}
	case constants.EndpointFs:
		dst, err = fs.New(ctx, constants.DestinationEndpoint)
		if err != nil {
			return
		}
	case constants.EndpointS3:
		dst, err = s3.New(ctx, constants.DestinationEndpoint, contexts.Client)
		if err != nil {
			return
		}
	default:
		logrus.Errorf("Type %s is not supported.", t.Src.Type)
		err = constants.ErrEndpointNotSupported
		return
	}

	return
}

// run will execute task.
func run(ctx context.Context) (err error) {
	// Check if task has been finished.
	if t.Status == constants.TaskStatusFinished {
		logrus.Infof("Task %s has been finished, skip.", t.Name)
		return
	}

	switch t.Type {
	case constants.TaskTypeCopy:
		t.Handle = copyObject
		err = copyTask(ctx)
		if err != nil {
			return
		}
	case constants.TaskTypeDelete:
		t.Handle = deleteObject
		err = deleteTask(ctx)
		if err != nil {
			return
		}
	case constants.TaskTypeFetch:
		t.Handle = fetchObject
		err = fetchTask(ctx)
		if err != nil {
			return
		}
	default:
		logrus.Errorf("Task %s's type %s is not supported.", t.Name, t.Type)
		return
	}

	// Update task status.
	t.Status = constants.TaskStatusFinished
	err = t.Save(ctx)
	if err != nil {
		logrus.Errorf("Task %s save failed for %v.", t.Name, err)
		return
	}

	logrus.Infof("Task %s has been finished.", t.Name)
	return
}

// migrateWorker will only do migrate work.
func migrateWorker(ctx context.Context) {
	defer owg.Done()
	defer utils.Recover()

	for o := range oc {
		ok, err := checkObject(ctx, o)
		if err != nil {
			logrus.Errorf("Check object failed for %v.", err)
			continue
		}
		if ok {
			err = model.DeleteObject(ctx, o)
			if err != nil {
				utils.CheckClosedDB(err)
			}
			continue
		}

		// Object may be tried in three times.
		bo := backoff.NewExponentialBackOff()
		bo.Multiplier = 2.0
		bo.MaxElapsedTime = 2 * time.Second

		err = backoff.Retry(func() error {
			err = t.Handle(ctx, o)
			if err == nil {
				return nil
			}

			logrus.Infof("%s object failed for %v, retried.", t.Type, err)
			return err
		}, bo)
		if err != nil {
			logrus.Errorf("%s object failed for %v.", t.Type, err)
			continue
		}

		err = model.DeleteObject(ctx, o)
		if err != nil {
			utils.CheckClosedDB(err)
			continue
		}
	}
}

// isFinished will check whether current task has been finished.
func isFinished(ctx context.Context) bool {
	h, err := model.HasDirectoryObject(ctx)
	if err != nil {
		logrus.Panic(err)
	}
	if h {
		logrus.Infof("There are not finished directory objects.")
		return false
	}

	h, err = model.HasSingleObject(ctx)
	if err != nil {
		logrus.Panic(err)
	}
	if h {
		logrus.Infof("There are not finished single objects.")
		return false
	}

	h, err = model.HasPartialObject(ctx)
	if err != nil {
		logrus.Panic(err)
	}
	if h {
		logrus.Infof("There are not finished partial objects.")
		return false
	}

	return true
}
