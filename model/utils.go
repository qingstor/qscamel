package model

import (
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/db"
	"github.com/yunify/qscamel/utils"
)

// CloseTx will close tx with err.
// If err is nil, we will try to commit this tx.
// If err is not nil, we will rollback.
func CloseTx(tx *db.Tx, err error) {
	defer utils.Recover()

	// If not writable, just rollback and skip.
	if !tx.Writable() {
		tx.Rollback()
		return
	}

	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logrus.Errorf("Tx failed to commit for %v.", err)
	}
}
