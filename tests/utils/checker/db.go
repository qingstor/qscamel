package checker

import (
	"fmt"

	"github.com/yunify/qscamel/db"
)

// CheckDBEmpty check the temp dbfile and testing
// whether the Database is empty.
func CheckDBEmpty(baseDir string) error {
	db, err := db.NewDB(&db.DatabaseOptions{baseDir + "/db"})
	if err != nil {
		return err
	}
	defer db.Close()
	it := db.NewIterator(nil, nil)
	if it.Next() == true {
		return dbtest{"database is not empty"}
	}
	return nil
}

type dbtest struct {
	fail string
}

func (e dbtest) Error() string {
	return fmt.Sprintf("%s", e.fail)
}
