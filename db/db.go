package db

import (
	"time"

	"github.com/coreos/bbolt"
	"github.com/sirupsen/logrus"

	"github.com/yunify/qscamel/constants"
)

// Database stores database connection.
type Database struct {
	*bolt.DB
}

// DatabaseOptions stores database options.
type DatabaseOptions struct {
	Address string
}

// NewDB will create a new database connection.
func NewDB(opt *DatabaseOptions) (d *Database, err error) {
	client, err := bolt.Open(opt.Address, 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logrus.Errorf("Open database failed for %v.", err)
		return
	}

	d = &Database{DB: client}
	logrus.Debugf("Connected to Bolt database %s", opt.Address)
	return
}

// Init will init current database.
func (db *Database) Init() (err error) {
	tx, err := db.Begin(true)
	if err != nil {
		logrus.Errorf("Start writable transaction failed for %v.", err)
		return
	}
	defer tx.Rollback()

	// Create task list bucket.
	_, err = tx.CreateBucketIfNotExists([]byte(constants.KeyTaskList))
	if err != nil {
		logrus.Errorf("Create task bucket failed for %v.", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		logrus.Errorf("DB commit failed for %v.", err)
		return
	}

	return
}
