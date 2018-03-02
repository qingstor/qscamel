package model

import (
	"context"
	"strconv"

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
	Status string `msgpack:"s"`
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
		ID:   id,
		Path: p,
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

// GetCurrentJobID will get current job ID.
func GetCurrentJobID(ctx context.Context) (id uint64, err error) {
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

	val := b.Get([]byte(constants.KeyCurrentJob))
	if val == nil {
		return 0, nil
	}
	id, err = strconv.ParseUint(string(val), 10, 64)
	if err != nil {
		logrus.Panic("ParseUint failed for %v.", err)
	}
	return
}

// UpdateCurrentJobID will update current job ID.
func UpdateCurrentJobID(ctx context.Context, n uint64) (err error) {
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

	return b.Put([]byte(constants.KeyCurrentJob), []byte(strconv.FormatUint(n, 10)))
}
