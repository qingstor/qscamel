package model

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/contexts"
	"github.com/yunify/qscamel/utils"
)

// GetSequence will get current bucket's sequence.
func GetSequence(ctx context.Context) (n uint64, err error) {
	t := utils.FromTaskContext(ctx)

	tx := utils.FromTxContext(ctx)
	if tx == nil {
		tx, err = contexts.DB.Begin(false)
		if err != nil {
			logrus.Errorf("Start writable transaction failed for %v.", err)
			return
		}
		defer func() {
			CloseTx(tx, err)
		}()
	}

	b := tx.Bucket(constants.FormatTaskKey(t))
	n = b.Sequence()

	return
}

// NextSequence will return bucket's next sequence.
func NextSequence(ctx context.Context) (n uint64, err error) {
	t := utils.FromTaskContext(ctx)

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

	b := tx.Bucket(constants.FormatTaskKey(t))
	n, err = b.NextSequence()
	if err != nil {
		return
	}

	return
}
