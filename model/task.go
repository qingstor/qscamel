package model

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/utils"
)

// Task store all data for a task.
type Task struct {
	Type string `yaml:"type" msgpack:"t"`

	Src *Endpoint `yaml:"source" msgpack:"src"`
	Dst *Endpoint `yaml:"destination" msgpack:"dst"`

	IgnoreExisting        string `yaml:"ignore_existing" msgpack:"ie"`
	MultipartBoundarySize int64  `yaml:"multipart_boundary_size" msgpack:"mbs"`
	// Format: 2006-01-02 15:04:05
	IgnoreBefore          string `yaml:"ignore_before" msgpack:"ib"`
	IgnoreBeforeTimestamp int64  `yaml:"-" msgpack:"ibt"`
	RateLimit             int    `yaml:"rate_limit" msgpack:"rl"`

	// Statistical Information
	SuccessCount  int64               `yaml:"-" msgpack:"sc"`
	SuccessSize   int64               `yaml:"-" msgpack:"ss"`
	FailedObjects map[string]struct{} `yaml:"-" msgpack:"fo"`

	// Data that only stores in database.
	Name   string `yaml:"-" msgpack:"n"`
	Status string `yaml:"-" msgpack:"s"`

	// Date that only keep in memory.
	Handle func(ctx context.Context, o Object) (err error) `yaml:"-" msgpack:"-"`
}

// LoadTask will try to load task from database and file.
func LoadTask(name, taskPath string) (t *Task, err error) {
	// Load from database first.
	t, err = GetTaskByName(nil, name)
	if err != nil {
		return
	}

	if taskPath == "" {
		if t == nil {
			// If t is nil and no task path input, we should return not found error.
			return nil, constants.ErrTaskNotFound
		}
		// If t is found and no task path input, we should return the task.
		return t, nil
	}

	// Load from file
	task, err := LoadTaskFromFilePath(taskPath)
	if err != nil {
		return
	}

	// If t is not nil and task path input, we should check the task content.
	if t != nil {
		if t.Sum256() != task.Sum256() {
			return nil, constants.ErrTaskMismatch
		}
		return t, nil
	}

	// Parse ignore before
	if task.IgnoreBefore != "" {
		format := "2006-01-02 15:04:05"
		ignoreBefore, err := time.Parse(format, task.IgnoreBefore)
		if err != nil {
			logrus.Errorf("%s is not a valid value for task ignore before", task.IgnoreBefore)
			return nil, constants.ErrTaskInvalid
		}
		task.IgnoreBeforeTimestamp = ignoreBefore.Unix()

	}

	if task.RateLimit == 0 {
		task.RateLimit = 1000
	}

	// Init FailedObjects map
	task.FailedObjects = make(map[string]struct{})

	// If task not in database, set task status to
	// created and save it.
	task.Status = constants.TaskStatusCreated
	task.Name = name
	err = task.Save(nil)
	if err != nil {
		return
	}
	return task, nil
}

// LoadTaskFromFilePath will load config from specific file path.
func LoadTaskFromFilePath(filePath string) (t *Task, err error) {
	filePath, err = utils.Expand(filePath)
	if err != nil {
		return
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	return LoadTaskFromContent(content)
}

// LoadTaskFromContent will load config from file content.
func LoadTaskFromContent(content []byte) (t *Task, err error) {
	t = &Task{}
	err = yaml.Unmarshal(content, t)
	if err != nil {
		return
	}
	return
}

// Check will check whether current task is valid.
func (t *Task) Check() error {
	switch t.IgnoreExisting {
	case "":
	case constants.TaskIgnoreExistingLastModified:
	case constants.TaskIgnoreExistingMD5Sum:
	default:
		logrus.Errorf("%s is not a valid value for task ignore existing", t.IgnoreExisting)
		return constants.ErrTaskInvalid
	}

	return nil
}

// Sum256 will calculate task's sha256.
func (t *Task) Sum256() [sha256.Size]byte {
	y, err := yaml.Marshal(t)
	if err != nil {
		logrus.Panicf("YAML marshal failed for %v.", err)
	}
	return sha256.Sum256(y)
}

// Save will save current task in DB.
func (t *Task) Save(ctx context.Context) (err error) {
	content, err := msgpack.Marshal(t)
	if err != nil {
		logrus.Panicf("Msgpack marshal failed for %v.", err)
	}

	err = contexts.DB.Put(constants.FormatTaskKey(t.Name), content, nil)
	if err != nil {
		return
	}

	return
}

// GetTask will get task by it's name.
func GetTask(ctx context.Context) (t *Task, err error) {
	name := utils.FromTaskContext(ctx)
	return GetTaskByName(ctx, name)
}

// GetTaskByName will get task by it's name.
func GetTaskByName(ctx context.Context, p string) (t *Task, err error) {
	t = &Task{}

	content, err := contexts.DB.Get(constants.FormatTaskKey(p), nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}

	err = msgpack.Unmarshal(content, t)
	if err != nil {
		logrus.Errorf("Msgpack marshal task %s failed for %v.", p, err)
		return
	}
	return
}

// DeleteTask will delete a task from content.
func DeleteTask(ctx context.Context) (err error) {
	name := utils.FromTaskContext(ctx)
	return DeleteTaskByName(ctx, name)
}

// DeleteTaskByName will delete a task by it's name.
func DeleteTaskByName(ctx context.Context, p string) (err error) {
	x := ""
	for {
		j, err := NextDirectoryObject(ctx, x)
		if err != nil {
			return err
		}
		if j == nil {
			break
		}

		err = DeleteObject(ctx, j)
		if err != nil {
			return err
		}

		x = j.Key

		logrus.Infof("Task %s, directory object %s has been deleted.", p, j.Key)
	}

	x = ""
	for {
		o, err := NextSingleObject(ctx, x)
		if err != nil {
			return err
		}
		if o == nil {
			break
		}

		err = DeleteObject(ctx, o)
		if err != nil {
			return err
		}

		x = o.Key

		logrus.Infof("Task %s, single object %s has been deleted.", p, o.Key)
	}

	x = ""
	pn := -1
	for {
		po, err := NextPartialObject(ctx, x, pn)
		if err != nil {
			logrus.Panic(err)
		}
		if po == nil {
			break
		}

		err = DeleteObject(ctx, po)
		if err != nil {
			return err
		}

		x = po.Key
		pn = po.PartNumber

		logrus.Infof("Task %s, partial object %s at %d has been deleted.",
			p, po.Key, po.PartNumber)
	}

	err = contexts.DB.Delete(constants.FormatTaskKey(p), nil)
	if err != nil {
		return
	}
	return
}

// ListTask will list all tasks.
func ListTask(ctx context.Context) (t []*Task, err error) {
	t = []*Task{}

	it := contexts.DB.NewIterator(
		util.BytesPrefix(constants.FormatTaskKey("")), nil)

	for it.Next() {
		k := it.Key()

		// If k doesn't has object prefix, there are no job any more.
		if !bytes.HasPrefix(k, constants.FormatTaskKey("")) {
			break
		}

		task := &Task{}

		v := it.Value()
		err = msgpack.Unmarshal(v, task)
		if err != nil {
			logrus.Panicf("Msgpack unmarshal failed for %v.", err)
		}

		t = append(t, task)
	}

	it.Release()
	err = it.Error()
	return
}
