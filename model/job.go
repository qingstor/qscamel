package model

import (
	"bytes"
	"context"

	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/vmihailenco/msgpack"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/utils"
)

// Job stores job status.
type Job struct {
	Path   string `msgpack:"p"`
	Marker string `msgpack:"m"`
}

// Save will save current job to DB.
func (j *Job) Save(ctx context.Context) (err error) {
	t := utils.FromTaskContext(ctx)

	content, err := msgpack.Marshal(j)
	if err != nil {
		return
	}

	err = contexts.DB.Put(constants.FormatJobKey(t, j.Path), content, nil)
	if err != nil {
		return
	}

	return
}

// CreateJob will create a new job.
func CreateJob(ctx context.Context, p string) (j *Job, err error) {
	t := utils.FromTaskContext(ctx)

	j = &Job{
		Path:   p,
		Marker: "",
	}

	content, err := msgpack.Marshal(j)
	if err != nil {
		return
	}

	err = contexts.DB.Put(constants.FormatJobKey(t, j.Path), content, nil)
	if err != nil {
		return
	}

	return
}

// DeleteJob will delete a job.
func DeleteJob(ctx context.Context, p string) (err error) {
	t := utils.FromTaskContext(ctx)

	return contexts.DB.Delete(constants.FormatJobKey(t, p), nil)
}

// GetJob will get a job by it's ID.
func GetJob(ctx context.Context, p string) (j *Job, err error) {
	t := utils.FromTaskContext(ctx)

	j = &Job{}

	content, err := contexts.DB.Get(constants.FormatJobKey(t, p), nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return
	}

	err = msgpack.Unmarshal(content, j)
	if err != nil {
		logrus.Panicf("Msgpack unmarshal failed for %v.", err)
	}
	return
}

// HasJob will check whether db has not finished job.
func HasJob(ctx context.Context) (b bool, err error) {
	t := utils.FromTaskContext(ctx)

	it := contexts.DB.NewIterator(
		util.BytesPrefix(constants.FormatJobKey(t, "")), nil)

	b = it.Seek(constants.FormatJobKey(t, ""))

	if b {
		key := it.Key()

		if !bytes.HasPrefix(key, constants.FormatJobKey(t, "")) {
			b = false
		}
	}

	it.Release()
	err = it.Error()
	return
}

// NextJob will return the next job after id.
func NextJob(ctx context.Context, p string) (j *Job, err error) {
	t := utils.FromTaskContext(ctx)

	it := contexts.DB.NewIterator(
		util.BytesPrefix(constants.FormatJobKey(t, "")), nil)

	for ok := it.Seek(constants.FormatJobKey(t, p)); ok; ok = it.Next() {
		k := it.Key()

		// Check if the same key first, and go further.
		if bytes.Compare(k, constants.FormatJobKey(t, p)) == 0 {
			continue
		}
		// If k doesn't has job prefix, there are no job any more.
		if !bytes.HasPrefix(k, constants.FormatJobKey(t, "")) {
			break
		}

		j = &Job{}
		v := it.Value()
		err = msgpack.Unmarshal(v, j)
		if err != nil {
			logrus.Panicf("Msgpack unmarshal failed for %v.", err)
		}
		return
	}

	it.Release()
	err = it.Error()
	return
}
