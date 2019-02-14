package model

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/vmihailenco/msgpack"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/utils"
)

// Object is the interface for object.
type Object interface {
	Type() string
}

// DirectoryObject is the object for directory.
type DirectoryObject struct {
	Key    string `msgpack:"p"`
	Marker string `msgpack:"m"`
}

// Type implement Object.Type
func (o *DirectoryObject) Type() string {
	return constants.ObjectTypeDirectory
}

// SingleObject is a single object.
type SingleObject struct {
	Key string `msgpack:"p"`

	Size         int64  `msgpack:"s"`
	LastModified int64  `msgpack:"lm"`
	MD5          string `msgpack:"cm"`
}

// Type implement Object.Type
func (o *SingleObject) Type() string {
	return constants.ObjectTypeSingle
}

// PartialObject is a partial object.
type PartialObject struct {
	Key string `msgpack:"p"`

	Size   int64  `msgpack:"s"`
	Offset int64  `msgpack:"of"`
	MD5    string `msgpack:"cm"`

	TotalNumber int    `msgpack:"tn"`
	PartNumber  int    `msgpack:"pn"`
	UploadID    string `msgpack:"uid"`
}

// Type implement Object.Type
func (o *PartialObject) Type() string {
	return constants.ObjectTypePartial
}

// CreateObject will create an object in db.
func CreateObject(ctx context.Context, o Object) (err error) {
	t := utils.FromTaskContext(ctx)

	content, err := msgpack.Marshal(o)
	if err != nil {
		logrus.Panicf("Msgpack marshal failed for %v.", err)
	}

	switch x := o.(type) {
	case *DirectoryObject:
		return contexts.DB.Put(constants.FormatDirectoryObjectKey(t, x.Key), content, nil)
	case *SingleObject:
		return contexts.DB.Put(constants.FormatSingleObjectKey(t, x.Key), content, nil)
	case *PartialObject:
		return contexts.DB.Put(constants.FormatPartialObjectKey(t, x.Key, x.PartNumber), content, nil)
	default:
		err = constants.ErrObjectInvalid
		return
	}
}

// DeleteObject will delete an object.
func DeleteObject(ctx context.Context, o Object) (err error) {
	t := utils.FromTaskContext(ctx)

	switch x := o.(type) {
	case *DirectoryObject:
		return contexts.DB.Delete(constants.FormatDirectoryObjectKey(t, x.Key), nil)
	case *SingleObject:
		return contexts.DB.Delete(constants.FormatSingleObjectKey(t, x.Key), nil)
	case *PartialObject:
		return contexts.DB.Delete(constants.FormatPartialObjectKey(t, x.Key, x.PartNumber), nil)
	default:
		err = constants.ErrObjectInvalid
		return
	}
}

// HasDirectoryObject will check whether db has not finished directory object.
func HasDirectoryObject(ctx context.Context) (b bool, err error) {
	t := utils.FromTaskContext(ctx)
	return hasObject(ctx, constants.FormatDirectoryObjectKey(t, ""))
}

// HasSingleObject will check whether db has not finished single object.
func HasSingleObject(ctx context.Context) (b bool, err error) {
	t := utils.FromTaskContext(ctx)
	return hasObject(ctx, constants.FormatSingleObjectKey(t, ""))
}

// HasPartialObject will check whether db has not finished partial object.
func HasPartialObject(ctx context.Context) (b bool, err error) {
	t := utils.FromTaskContext(ctx)
	return hasObject(ctx, constants.FormatPartialObjectKey(t, "", -1))
}

// HasParts will check whether specific key has not finished parts.
func HasParts(ctx context.Context, key string) (b bool, err error) {
	t := utils.FromTaskContext(ctx)
	return hasObject(ctx, constants.FormatPartialObjectKey(t, key, -1))
}

func hasObject(ctx context.Context, v []byte) (b bool, err error) {
	it := contexts.DB.NewIterator(
		util.BytesPrefix(v), nil)

	b = it.Seek(v)
	if b {
		key := it.Key()

		if !bytes.HasPrefix(key, v) {
			b = false
		}
	}

	it.Release()
	err = it.Error()
	return
}

// NextDirectoryObject will return the next directory object after p.
func NextDirectoryObject(ctx context.Context, p string) (o *DirectoryObject, err error) {
	t := utils.FromTaskContext(ctx)

	it := contexts.DB.NewIterator(
		util.BytesPrefix(constants.FormatDirectoryObjectKey(t, "")), nil)
	for ok := it.Seek(constants.FormatDirectoryObjectKey(t, p)); ok; ok = it.Next() {
		k := it.Key()

		// Check if the same key first, and go further.
		if bytes.Compare(k, constants.FormatDirectoryObjectKey(t, p)) == 0 {
			continue
		}
		// If k doesn't has object prefix, there are no object any more.
		if !bytes.HasPrefix(k, constants.FormatDirectoryObjectKey(t, "")) {
			break
		}

		o = &DirectoryObject{}
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

// NextSingleObject will return the next single object after p.
func NextSingleObject(ctx context.Context, p string) (o *SingleObject, err error) {
	t := utils.FromTaskContext(ctx)

	it := contexts.DB.NewIterator(
		util.BytesPrefix(constants.FormatSingleObjectKey(t, "")), nil)

	for ok := it.Seek(constants.FormatSingleObjectKey(t, p)); ok; ok = it.Next() {
		k := it.Key()

		// Check if the same key first, and go further.
		if bytes.Compare(k, constants.FormatSingleObjectKey(t, p)) == 0 {
			continue
		}
		// If k doesn't has object prefix, there are no object any more.
		if !bytes.HasPrefix(k, constants.FormatSingleObjectKey(t, "")) {
			break
		}

		o = &SingleObject{}
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

// NextPartialObject will return the next single object after p.
func NextPartialObject(ctx context.Context, p string, partNumber int) (o *PartialObject, err error) {
	t := utils.FromTaskContext(ctx)

	it := contexts.DB.NewIterator(
		util.BytesPrefix(constants.FormatPartialObjectKey(t, p, -1)), nil)

	for ok := it.Seek(constants.FormatPartialObjectKey(t, p, partNumber)); ok; ok = it.Next() {
		k := it.Key()

		// Check if the same key first, and go further.
		if bytes.Compare(k, constants.FormatPartialObjectKey(t, p, partNumber)) == 0 {
			continue
		}
		// If k doesn't has object prefix, there are no object any more.
		if !bytes.HasPrefix(k, constants.FormatPartialObjectKey(t, p, -1)) {
			break
		}

		o = &PartialObject{}
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
