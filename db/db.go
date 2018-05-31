package db

import (
	"github.com/sirupsen/logrus"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

// Database stores database connection.
type Database struct {
	*leveldb.DB
}

// DatabaseOptions stores database options.
type DatabaseOptions struct {
	Address string
}

// NewDB will create a new database connection.
func NewDB(opt *DatabaseOptions) (d *Database, err error) {
	// Set NoFreelistSync to true to import write performance.
	client, err := leveldb.OpenFile(opt.Address, nil)
	if err == nil {
		d = &Database{DB: client}
		logrus.Debugf("Connected to database %s", opt.Address)
		return
	}

	if _, ok := err.(*storage.ErrCorrupted); ok {
		logrus.Errorf("Open database failed for %v, recovering.", err)
		client, err = leveldb.RecoverFile(opt.Address, nil)
		if err != nil {
			logrus.Errorf("Database is corrupted and recover failed for %v.", err)
			return
		}
	}

	logrus.Errorf("Open database failed for %v.", err)
	return
}
