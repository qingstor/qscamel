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
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"go.uber.org/ratelimit"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/endpoint"
	"github.com/yunify/qscamel/endpoint/aliyun"
	"github.com/yunify/qscamel/endpoint/azblob"
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

	rl ratelimit.Limiter

	pool *ants.Pool

	multipartBoundarySize int64
)

// Execute will execute migrate task.
func Execute(ctx context.Context, close chan struct{}) (err error) {
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

	rl = ratelimit.New(t.RateLimit)

	var workers int
	if t.Workers == 0 {
		workers = 100
	} else {
		workers = t.Workers
	}
	pool, err = ants.NewPool(workers)
	if err != nil {
		logrus.Errorf("New migrate multipart workers failed for %v.", err)
		return
	}

	err = check(ctx)
	if err != nil {
		logrus.Errorf("Pre migrate check failed for %v.", err)
		return
	}

	return run(ctx, close)
}

func check(ctx context.Context) (err error) {
	// Initialize source.
	switch t.Src.Type {
	case constants.EndpointAliyun:
		src, err = aliyun.New(ctx, constants.SourceEndpoint, contexts.Client)
		if err != nil {
			return
		}
	case constants.EndpointAzblob:
		src, err = azblob.New(ctx, constants.SourceEndpoint, contexts.Client)
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
		logrus.Errorf("Type src %s is not supported.", t.Src.Type)
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
		logrus.Errorf("Type dst %s is not supported.", t.Dst.Type)
		err = constants.ErrEndpointNotSupported
		return
	}

	return
}

// run will execute task.
func run(ctx context.Context, close chan struct{}) (err error) {
	// Check if task has been finished.
	if t.Status == constants.TaskStatusFinished {
		logrus.Infof("Task %s has been finished, skip.", t.Name)
		return
	}

	go printStatistics(close)

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

	close <- struct{}{}

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
		backOff := backoff.WithMaxTries(bo, 10)

		fn := func() error {
			rl.Take()

			err = t.Handle(ctx, o)
			if err == nil {
				return nil
			}

			switch x := o.(type) {
			case *model.SingleObject:
				t.FailedObjects[x.Key] = 0
			}

			logrus.Infof("%s object failed for %v, retried.", t.Type, err)
			return err
		}

		err = backoff.Retry(fn, backOff)
		if err != nil {
			switch o.(type) {
			case *model.SingleObject:
				e := model.DeleteObject(ctx, o)
				if e != nil {
					utils.CheckClosedDB(e)
					continue
				}
			}
			logrus.Errorf("%s object failed for %v.", t.Type, err)
			continue
		}

		err = model.DeleteObject(ctx, o)
		if err != nil {
			utils.CheckClosedDB(err)
			continue
		}

		switch x := o.(type) {
		case *model.SingleObject:
			if _, ok := t.FailedObjects[x.Key]; ok {
				delete(t.FailedObjects, x.Key)
			}
			t.SuccessCount++
			t.SuccessSize += x.Size
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

func printStatistics(close chan struct{}) {
	timer := time.NewTicker(5 * time.Second)
	var tmpCount int64
	for {
		select {
		case <-close:
			logrus.Infof("====Final Success Count: %d  Final Success Size: %d====", t.SuccessCount, t.SuccessSize)
			filenames := make([]string, 0)
			for name, _ := range t.FailedObjects {
				filenames = append(filenames, name)
			}
			if len(filenames) > 0 {
				logrus.Infof("====Final Failed Count: %d  Final Failed filename: %v====", len(filenames), filenames)
			} else {
				logrus.Infof("====All objects migrated successfully====")
			}
			break
		case <-timer.C:
			if tmpCount != t.SuccessCount {
				logrus.Infof("====Success Count: %d  Success Size: %d====", t.SuccessCount, t.SuccessSize)
				tmpCount = t.SuccessCount
			}
		}
	}
}

func SaveTask() {
	if t.SuccessCount != 0 || len(t.FailedObjects) > 0 {
		_ = t.Save(nil)
	}
}
