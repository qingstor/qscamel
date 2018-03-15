package model

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/utils"
)

// Task store all data for a task.
type Task struct {
	Name string `yaml:"name" msgpack:"n"`
	Type string `yaml:"type" msgpack:"t"`

	Src *Endpoint `yaml:"source" msgpack:"src"`
	Dst *Endpoint `yaml:"destination" msgpack:"dst"`

	Overwrite        bool `yaml:"overwrite" msgpack:"o"`
	IgnoreExisting   bool `yaml:"ignore_existing" msgpack:"ie"`
	IgnoreUnmodified bool `yaml:"ignore_unmodified" msgpack:"iu"`

	// Data that only stores in database.
	Status string `yaml:"-" msgpack:"s"`
}

// LoadTask will try to load task from database and file.
func LoadTask(s string) (t *Task, err error) {
	// Load from database first.
	t, err = GetTaskByName(nil, s)
	if err != nil {
		return
	}
	if t != nil {
		return
	}

	// Load from file
	task, err := LoadTaskFromFilePath(s)
	if err != nil {
		return
	}
	t, err = GetTaskByName(nil, task.Name)
	if err != nil {
		return
	}
	if t == nil {
		// If task not in database, set task status to
		// running and save it.
		task.Status = constants.TaskStatusRunning
		err = task.Save(nil)
		if err != nil {
			return
		}
		return task, err
	}
	return
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
	// TODO: check other value.

	td, err := GetTaskByName(nil, t.Name)
	if err != nil {
		return err
	}
	if td != nil && t.Sum256() != td.Sum256() {
		logrus.Infof("Task content has been changed, check failed.")
		return constants.ErrTaskMismatch
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
	tx := utils.FromTxContext(ctx)
	if tx == nil {
		tx, err = contexts.DB.Begin(true)
		if err != nil {
			logrus.Errorf("Start writable transaction failed for %v.", err)
			return
		}
		defer func() {
			CloseTx(tx, err)
		}()
	}

	b := tx.Bucket([]byte(constants.KeyTaskList))

	content, err := msgpack.Marshal(t)
	if err != nil {
		logrus.Panicf("Msgpack marshal failed for %v.", err)
	}

	err = b.Put(constants.FormatTaskKey(t.Name), content)
	if err != nil {
		return
	}

	// Create related task bucket.
	_, err = tx.CreateBucketIfNotExists(constants.FormatTaskKey(t.Name))
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

	tx := utils.FromTxContext(ctx)
	if tx == nil {
		tx, err = contexts.DB.Begin(false)
		if err != nil {
			logrus.Errorf("Start read only transaction failed for %v.", err)
			return
		}
		defer func() {
			CloseTx(tx, err)
		}()
	}

	b := tx.Bucket([]byte(constants.KeyTaskList))

	content := b.Get(constants.FormatTaskKey(p))
	if content == nil {
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
	tx := utils.FromTxContext(ctx)
	if tx == nil {
		tx, err = contexts.DB.Begin(true)
		if err != nil {
			logrus.Errorf("Start read only transaction failed for %v.", err)
			return
		}
		defer func() {
			CloseTx(tx, err)
		}()
	}

	b := tx.Bucket([]byte(constants.KeyTaskList))
	err = b.Delete(constants.FormatTaskKey(p))
	if err != nil {
		return
	}

	return tx.DeleteBucket(constants.FormatTaskKey(p))
}

// ListTask will list all tasks.
func ListTask(ctx context.Context) (t []*Task, err error) {
	tx := utils.FromTxContext(ctx)
	if tx == nil {
		tx, err = contexts.DB.Begin(false)
		if err != nil {
			logrus.Errorf("Start transaction failed for %v.", err)
			return
		}
		defer func() {
			CloseTx(tx, err)
		}()
	}

	t = []*Task{}
	c := tx.Bucket([]byte(constants.KeyTaskList)).Cursor()

	k, v := c.Seek([]byte(constants.KeyTaskPrefix))
	for k != nil && bytes.HasPrefix(k, []byte(constants.KeyTaskPrefix)) {
		task := &Task{}
		err = msgpack.Unmarshal(v, task)
		if err != nil {
			logrus.Panicf("Msgpack unmarshal failed for %v.", err)
		}

		t = append(t, task)

		k, v = c.Next()
	}

	return
}
