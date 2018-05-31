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

// Object stores object info.
type Object struct {
	Key string `msgpack:"p"`

	IsDir        bool   `msgpack:"id"`
	Size         int64  `msgpack:"s"`
	LastModified int64  `msgpack:"lm"`
	MD5          string `msgpack:"cm"`
}

// CreateObject will create an object in db.
func CreateObject(ctx context.Context, o *Object) (err error) {
	t := utils.FromTaskContext(ctx)

	content, err := msgpack.Marshal(o)
	if err != nil {
		logrus.Panicf("Msgpack marshal failed for %v.", err)
	}

	return contexts.DB.Put(constants.FormatObjectKey(t, o.Key), content, nil)
}

// DeleteObject will delete an object.
func DeleteObject(ctx context.Context, p string) (err error) {
	t := utils.FromTaskContext(ctx)

	return contexts.DB.Delete(constants.FormatObjectKey(t, p), nil)
}

// GetObject will get an object from db.
func GetObject(ctx context.Context, p string) (o *Object, err error) {
	t := utils.FromTaskContext(ctx)

	o = &Object{}

	content, err := contexts.DB.Get(constants.FormatObjectKey(t, p), nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
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

	it := contexts.DB.NewIterator(
		util.BytesPrefix(constants.FormatObjectKey(t, "")), nil)

	b = it.Seek(constants.FormatObjectKey(t, ""))

	if b {
		key := it.Key()

		if !bytes.HasPrefix(key, constants.FormatObjectKey(t, "")) {
			b = false
		}
	}

	it.Release()
	err = it.Error()
	return
}

// NextObject will return the next object after p.
func NextObject(ctx context.Context, p string) (o *Object, err error) {
	t := utils.FromTaskContext(ctx)

	it := contexts.DB.NewIterator(
		util.BytesPrefix(constants.FormatObjectKey(t, "")), nil)

	for ok := it.Seek(constants.FormatObjectKey(t, p)); ok; ok = it.Next() {
		k := it.Key()

		// Check if the same key first, and go further.
		if bytes.Compare(k, constants.FormatObjectKey(t, p)) == 0 {
			continue
		}
		// If k doesn't has object prefix, there are no job any more.
		if !bytes.HasPrefix(k, constants.FormatObjectKey(t, "")) {
			break
		}

		o = &Object{}
		v := it.Value()
		err = msgpack.Unmarshal(v, o)
		if err != nil {
			logrus.Panicf("Msgpack unmarshal failed for %v.", err)
		}
		return
	}

	it.Release()
	err = it.Error()
	return
}
