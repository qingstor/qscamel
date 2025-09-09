package migrate

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/endpoint"
	"github.com/yunify/qscamel/model"
	"github.com/yunify/qscamel/utils"
)

func listObject(ctx context.Context, j *model.DirectoryObject) (err error) {
	defer jwg.Done()
	defer utils.Recover()

	srcName := src.Name(ctx)
	dstName := dst.Name(ctx)
	err = src.List(ctx, j, func(o model.Object) {
		defer utils.Recover()

		switch x := o.(type) {
		case *model.DirectoryObject:
			err = model.CreateObject(ctx, x)
			if err != nil {
				utils.CheckClosedDB(err)
			}

			logrus.Debugf("Directory object %s created.", x.Key)
			return
		case *model.SingleObject:
			if x.IsDir &&
				(!strings.Contains(srcName, "qingstor") && !strings.Contains(srcName, "s3")) &&
				(!strings.Contains(dstName, "qingstor") && !strings.Contains(dstName, "s3")) {
				return
			}
			err = model.CreateObject(ctx, x)
			if err != nil {
				utils.CheckClosedDB(err)
			}

			oc <- o
			return
		}
	})
	if err != nil {
		logrus.Errorf("Src list failed for %v.", err)
		return
	}

	err = model.DeleteObject(ctx, j)
	if err != nil {
		utils.CheckClosedDB(err)
	}
	return
}

// checkObject will tell whether an object is ok.
func checkObject(ctx context.Context, mo model.Object) (ok bool, err error) {
	if (t.IgnoreExisting == "" && t.IgnoreBeforeTimestamp == 0) || mo.Type() == constants.ObjectTypePartial {
		return false, nil
	}

	o := mo.(*model.SingleObject)

	logrus.Infof("Start checking object %s.", o.Key)

	so, err := statObject(ctx, src, o, false)
	if err != nil {
		return
	}
	if so == nil {
		return true, nil
	}

	do, err := statObject(ctx, dst, o, false)
	if err != nil {
		return
	}
	// Check existence.
	if do == nil {
		return
	}
	// Check size.
	//if so.Size != do.Size {
	//	logrus.Infof("Object %s size is not match, execute an operation on it.", o.Key)
	//	return
	//}

	if t.IgnoreBeforeTimestamp != 0 {
		ignoreTs := t.IgnoreBeforeTimestamp

		switch t.IgnoreExisting {
		case constants.TaskIgnoreExistingLastModified:
			if so.LastModified > ignoreTs && so.LastModified > do.LastModified {
				logrus.Infof("Object %s was modified, execute an operation on it.", o.Key)
				return
			}

		case constants.TaskIgnoreExistingMD5Sum:
			if so.LastModified > ignoreTs && so.MD5 != do.MD5 {
				logrus.Infof("Object %s was modified, execute an operation on it.", o.Key)
				return
			}

		default:
			if so.LastModified > ignoreTs {
				logrus.Infof("Object %s was modified after %s, execute an operation on it.",
					o.Key, time.Unix(ignoreTs, 0))
				return
			}
		}

		logrus.Infof("Object %s check passed, ignore.", o.Key)
		return true, nil
	}

	// Check last modified
	if t.IgnoreExisting == constants.TaskIgnoreExistingLastModified {
		if so.LastModified > do.LastModified {
			logrus.Infof("Object %s was modified, execute an operation on it.", o.Key)
			return
		}
		logrus.Infof("Object %s check passed, ignore.", o.Key)
		return true, nil
	}

	// Check md5.
	if so.MD5 != do.MD5 {
		logrus.Infof("Object %s md5 is not match, execute an operation on it.", o.Key)
		return
	}

	logrus.Infof("Object %s check passed, ignore.", o.Key)
	return true, nil
}

// checkObjectAfterMigrate will check whether the MD5 between the migrated src and dst is consistent.
func checkObjectAfterMigrate(ctx context.Context, mo model.Object) (err error) {
	if t.Src.Type == "fs" || t.Dst.Type == "fs" {
		return nil
	}

	o := mo.(*model.SingleObject)

	rso, err := statObject(ctx, src, o, true)
	if err != nil {
		logrus.Errorf("Src stat %s failed for %v.", o.Key, err)
		return err
	}

	rdo, err := statObject(ctx, dst, o, true)
	if err != nil {
		logrus.Errorf("Dst stat %s failed for %v.", o.Key, err)
		return err
	}

	if rdo.MD5 != rso.MD5 {
		logrus.Errorf("md5 mismatch between src and dst %s.", o.Key)
		return fmt.Errorf("md5 not match")
	}

	return
}

// copyObject will do a real copy.
func copyObject(ctx context.Context, o model.Object) (err error) {
	so := o.(*model.SingleObject)

	logrus.Infof("Start copying object %s.", so.Key)

	// Upload single object, if don't to split it.
	if so.Size <= multipartBoundarySize || !dst.Partable() {
		r, err := src.Read(ctx, so.Key, so.IsDir)
		if err != nil {
			logrus.Errorf("Src read %s failed for %v.", so.Key, err)
			return err
		}
		err = dst.Write(ctx, so.Key, so.Size, r, so.IsDir, so.QSMetadata)
		if err != nil {
			logrus.Errorf("Dst write %s failed for %v.", so.Key, err)
			return err
		}

		if t.CheckMD5 {
			err = checkObjectAfterMigrate(ctx, o)
			if err != nil {
				_ = dst.Delete(ctx, so.Key)
				return err
			}
		}

		logrus.Infof("Single object %s copied.", so.Key)
		return nil
	}

	// Split single object into part objects.
	uploadID, partSize, partNumbers, err := dst.InitPart(ctx, so.Key, so.Size, so.QSMetadata)
	if err != nil {
		logrus.Errorf("Dst init part %s failed for %v.", so.Key, err)
		return err
	}

	var e error
	once := sync.Once{}
	// error exit
	eQuit := make(chan struct{})
	// Normal exit
	quit := make(chan struct{})
	go func() {
		defer close(quit)
		wg := sync.WaitGroup{}

		offset := int64(0)
		for i := 0; i < partNumbers; i++ {
			wg.Add(1)

			oo := &model.PartialObject{
				Key: so.Key,

				Size:   partSize,
				Offset: offset,

				TotalNumber: partNumbers,
				PartNumber:  i,
				UploadID:    uploadID,
			}

			offset += partSize
			if offset > so.Size {
				oo.Size = so.Size - offset + partSize
			}

			err = pool.Submit(func() {
				defer wg.Done()

				logrus.Infof("Start copying partial object %s at %d.", oo.Key, oo.PartNumber)

				r, err := src.ReadRange(ctx, oo.Key, oo.Offset, oo.Size)
				if err != nil {
					once.Do(func() {
						logrus.Errorf("Src read partial object %s at %d failed for %v.",
							oo.Key, oo.Offset, err)
						close(eQuit)
						e = err
					})
					return
				}
				err = dst.UploadPart(ctx, oo, r)
				if err != nil {
					once.Do(func() {
						logrus.Errorf("Dst write partial object %s at %d failed for %v.",
							oo.Key, oo.Offset, err)
						close(eQuit)
						e = err
					})
					return
				}

				logrus.Infof("Partial object %s at %d copied.", oo.Key, oo.PartNumber)
				return

			})
			if err != nil {
				once.Do(func() {
					logrus.Errorf("Submit Upload partial object %s request failed for %v", so.Key, err)
					close(eQuit)
					e = err
				})
				return
			}
		}

		wg.Wait()
	}()

	select {
	case <-eQuit:
		break
	case <-quit:
		break
	}

	if e != nil {
		err = dst.AbortUploads(ctx, so.Key, uploadID)
		if err != nil {
			logrus.Errorf("Abort partial object %s failed for %v", so.Key, err)
		}
		return e
	}

	err = dst.CompleteParts(ctx, so.Key, uploadID, partNumbers)
	if err != nil {
		logrus.Errorf("Complete partial object %s failed for %v", so.Key, err)
		return err
	}

	if t.CheckMD5 {
		if t.CheckMD5 {
			err = checkObjectAfterMigrate(ctx, o)
			if err != nil {
				_ = dst.Delete(ctx, so.Key)
				return err
			}
		}
	}

	logrus.Infof("Object %s copied.", so.Key)

	return
}

// deleteObject will do a real delete.
func deleteObject(ctx context.Context, o model.Object) (err error) {
	switch x := o.(type) {
	case *model.SingleObject:
		logrus.Infof("Start deleting single object %s.", x.Key)

		err = dst.Delete(ctx, x.Key)
		if err != nil {
			logrus.Errorf("Dst delete %s failed for %v.", x.Key, err)
			return err
		}

		logrus.Infof("Single object %s deleted.", x.Key)
	case *model.PartialObject:
		// TODO: we should handle delete partial object here.
	}

	return
}

// fetchObject will do a real fetch.
func fetchObject(ctx context.Context, o model.Object) (err error) {
	switch x := o.(type) {
	case *model.SingleObject:
		logrus.Infof("Start fetching single object %s.", x.Key)

		url, err := src.Reach(ctx, x.Key)
		if err != nil {
			logrus.Errorf("Src reach %s failed for %v.", x.Key, err)
			return err
		}
		err = dst.Fetch(ctx, x.Key, url)
		if err != nil {
			logrus.Errorf("Dst fetch %s failed for %v.", x.Key, err)
			return err
		}

		logrus.Infof("Single object %s fetched.", x.Key)
	case *model.PartialObject:
		logrus.Errorf("Object %s is invalid for fetch.", x.Key)
		err = constants.ErrObjectInvalid
		return
	}

	return
}

// statObject will get an object metadata and try to get it's md5 if available.
func statObject(
	ctx context.Context, e endpoint.Base, o *model.SingleObject, isMD5 bool,
) (ro *model.SingleObject, err error) {
	ro, err = e.Stat(ctx, o.Key, o.IsDir)
	if err != nil {
		logrus.Errorf("%s stat object %s failed for %v.", e.Name(ctx), o.Key, err)
		return
	}
	if ro == nil {
		logrus.Infof("Object %s is not found at %s.", o.Key, e.Name(ctx))
		return
	}

	if t.IgnoreExisting != constants.TaskIgnoreExistingMD5Sum && !isMD5 {
		return
	}

	if len(ro.MD5) != 32 {
		ro.MD5, err = md5SumObject(ctx, e, o)
		if err != nil {
			logrus.Errorf(
				"%s calculate object %s md5 failed for %v.", e.Name(ctx), o.Key, err)
			return
		}
	}
	return
}

// md5SumObject will get the object's md5
func md5SumObject(
	ctx context.Context, e endpoint.Base, o model.Object,
) (m string, err error) {
	var r io.Reader
	switch x := o.(type) {
	case *model.SingleObject:
		r, err = e.Read(ctx, x.Key, x.IsDir)
		if err != nil {
			logrus.Errorf("%s read single object %s failed for %v.",
				e.Name(ctx), x.Key, err)
			return
		}
	case *model.PartialObject:
		r, err = e.ReadRange(ctx, x.Key, x.Offset, x.Size)
		if err != nil {
			logrus.Errorf("%s read partial object %s aat %d failed for %v.",
				e.Name(ctx), x.Key, x.Offset, err)
			return
		}
	}

	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	sum := h.Sum(nil)

	return hex.EncodeToString(sum[:]), nil
}
