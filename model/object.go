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

// Object stores object info.
type Object struct {
	Key string `msgpack:"p"`

	IsDir        bool   `msgpack:"id"`
	Size         int64  `msgpack:"s"`
	LastModified int64  `msgpack:"lm"`
	ContentMD5   []byte `msgpack:"cm"`
}

// CreateObject will create an object in db.
func CreateObject(ctx context.Context, o *Object) (err error) {
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

	content, err := msgpack.Marshal(o)
	if err != nil {
		logrus.Panicf("Msgpack marshal failed for %v.", err)
	}

	return b.Put(constants.FormatObjectKey(o.Key), content)
}

// DeleteObject will delete an object.
func DeleteObject(ctx context.Context, p string) (err error) {
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

	return b.Delete(constants.FormatObjectKey(p))
}

// GetObject will get an object from db.
func GetObject(ctx context.Context, p string) (o *Object, err error) {
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

	o = &Object{}

	content := b.Get(constants.FormatObjectKey(p))
	if content == nil {
		return
	}

	err = msgpack.Unmarshal(content, o)
	if err != nil {
		logrus.Panicf("Msgpack unmarshal failed for %v.", err)
	}

	return
}

// HasObject will check whether db has not finished object.
func HasObject(ctx context.Context) (b bool, err error) {
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

	k, _ := c.Seek([]byte(constants.KeyObjectPrefix))

	if k != nil && bytes.HasPrefix(k, []byte(constants.KeyObjectPrefix)) {
		return true, nil
	}
	return
}

// ListObject will list current task's object.
func ListObject(ctx context.Context, fn func(*Object)) (err error) {
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
	o := &Object{}

	k, v := c.Seek([]byte(constants.KeyObjectPrefix))
	for k != nil && bytes.HasPrefix(k, []byte(constants.KeyObjectPrefix)) {
		err = msgpack.Unmarshal(v, o)
		if err != nil {
			logrus.Panicf("Msgpack unmarshal failed for %v.", err)
		}

		fn(o)

		k, v = c.Next()
	}

	return
}
