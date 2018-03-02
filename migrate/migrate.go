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
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/endpoint"
	"github.com/yunify/qscamel/endpoint/fs"
	"github.com/yunify/qscamel/endpoint/qingstor"
	"github.com/yunify/qscamel/model"
)

var (
	t *model.Task

	src endpoint.Source
	dst endpoint.Destination
)

// Execute will execute migrate task.
func Execute(ctx context.Context) (err error) {
	t, err = model.GetTask(ctx)
	if err != nil {
		return
	}

	if t.Status == constants.TaskStatusFinished {
		logrus.Infof("Task %s has been finished, skip.", t.Name)
		return
	}

	// Initialize source.
	switch t.Src.Type {
	case constants.EndpointQingStor:
		src, err = qingstor.New(ctx, constants.SourceEndpoint)
		if err != nil {
			return
		}
	case constants.EndpointFs:
		src, err = fs.New(ctx, constants.SourceEndpoint)
		if err != nil {
			return
		}
	default:
		logrus.Errorf("Type %s is not supported.", t.Src.Type)
		err = errors.New("type is not supported")
		return
	}

	// Initialize destination.
	switch t.Dst.Type {
	case constants.EndpointQingStor:
		dst, err = qingstor.New(ctx, constants.DestinationEndpoint)
		if err != nil {
			return
		}
	case constants.EndpointFs:
		dst, err = fs.New(ctx, constants.DestinationEndpoint)
		if err != nil {
			return
		}
	default:
		logrus.Errorf("Type %s is not supported.", t.Src.Type)
		err = errors.New("type is not supported")
		return
	}

	return Run(ctx)
}
