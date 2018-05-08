package model

import (
	"bytes"
	"context"

	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/utils"
)

// Job stores job status.
type Job struct {
	ID     uint64 `msgpack:"id"`
	Path   string `msgpack:"p"`
	Marker string `msgpack:"m"`
}

// Save will save current job to DB.
func (j *Job) Save(ctx context.Context) (err error) {
	t := utils.FromTaskContext(ctx)

	tx := utils.FromTxContext(ctx)
	if tx == nil {
		tx, err = contexts.DB.Begin(true)
		if err != nil {
			logrus.Errorf("Start transaction failed for %v.", err)
			return
		}
		defer func() {
			CloseTx(tx, err)
		}()

	}

	// Get bucket.
	b := tx.Bucket(constants.FormatTaskKey(t))

	content, err := msgpack.Marshal(j)
	if err != nil {
		return
	}

	err = b.Put(constants.FormatJobKey(j.ID), content)
	if err != nil {
		return
	}

	return
}

// CreateJob will create a new job.
func CreateJob(ctx context.Context, p string) (j *Job, err error) {
	t := utils.FromTaskContext(ctx)

	tx := utils.FromTxContext(ctx)
	if tx == nil {
		tx, err = contexts.DB.Begin(true)
		if err != nil {
			logrus.Errorf("Start transaction failed for %v.", err)
			return
		}
		defer func() {
			CloseTx(tx, err)
		}()
	}

	// Get bucket.
	b := tx.Bucket(constants.FormatTaskKey(t))

	id, err := b.NextSequence()
	if err != nil {
		return
	}

	j = &Job{
		ID:     id,
		Path:   p,
		Marker: "",
	}

	content, err := msgpack.Marshal(j)
	if err != nil {
		return
	}

	err = b.Put(constants.FormatJobKey(id), content)
	if err != nil {
		return
	}

	return
}

// DeleteJob will delete a job.
func DeleteJob(ctx context.Context, id uint64) (err error) {
	t := utils.FromTaskContext(ctx)

	tx := utils.FromTxContext(ctx)
	if tx == nil {
		tx, err = contexts.DB.Begin(true)
		if err != nil {
			logrus.Errorf("Start transaction failed for %v.", err)
			return
		}
		defer func() {
			CloseTx(tx, err)
		}()
	}

	b := tx.Bucket(constants.FormatTaskKey(t))

	return b.Delete(constants.FormatJobKey(id))
}

// GetJob will get a job by it's ID.
func GetJob(ctx context.Context, id uint64) (j *Job, err error) {
	t := utils.FromTaskContext(ctx)

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

	b := tx.Bucket(constants.FormatTaskKey(t))

	j = &Job{}

	err = msgpack.Unmarshal(b.Get(constants.FormatJobKey(id)), j)
	if err != nil {
		logrus.Panicf("Msgpack unmarshal failed for %v.", err)
	}

	return
}

// HasJob will check whether db has not finished job.
func HasJob(ctx context.Context) (b bool, err error) {
	t := utils.FromTaskContext(ctx)

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

	c := tx.Bucket(constants.FormatTaskKey(t)).Cursor()

	k, _ := c.Seek([]byte(constants.KeyJobPrefix))

	if k != nil && bytes.HasPrefix(k, []byte(constants.KeyJobPrefix)) {
		return true, nil
	}
	return
}

// NextJob will return the next job after id.
func NextJob(ctx context.Context, id uint64) (j *Job, err error) {
	t := utils.FromTaskContext(ctx)

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

	c := tx.Bucket(constants.FormatTaskKey(t)).Cursor()
	k, v := c.Seek(constants.FormatJobKey(id))

	// If k equal to current id, we should get the next id.
	if k != nil && bytes.Compare(k, constants.FormatJobKey(id)) == 0 {
		k, v = c.Next()
	}

	if k != nil && bytes.HasPrefix(k, []byte(constants.KeyJobPrefix)) {
		j = &Job{}
		err = msgpack.Unmarshal(v, j)
		if err != nil {
			logrus.Panicf("Msgpack unmarshal failed for %v.", err)
		}
		return
	}

	return
}
