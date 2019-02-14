package utils

import (
	"fmt"
	"testing"

	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/yunify/qscamel/db"
)

// CheckDBEmpty check the temp dbfile and testing
// whether the Database is empty.
func CheckDBEmpty(t *testing.T, fmap map[string]string) {
	resultDB, err := db.NewDB(&db.DatabaseOptions{Address: fmap["dir"] + "/db"})
	if err != nil {
		t.Fatal(err)
	}
	defer resultDB.Close()
	it := resultDB.NewIterator(nil, nil)
	if it.Next() == true {
		t.Fatal("database is not empty")
	}
	it.Release()
}

// CheckDBNoObject check the temp dbfile and testing
// whether the Database has any object
func CheckDBNoObject(t *testing.T, fmap map[string]string) {
	resultDB, err := db.NewDB(&db.DatabaseOptions{Address: fmap["dir"] + "/db"})
	if err != nil {
		t.Fatal(err)
	}
	defer resultDB.Close()

	taskName := fmap["name"]
	it := resultDB.NewIterator(util.BytesPrefix([]byte(fmt.Sprintf("~%s:do:", taskName))), nil)
	if it.Next() == true {
		t.Fatal("there still has directory object in database")
	}
	it.Release()

	it = resultDB.NewIterator(util.BytesPrefix([]byte(fmt.Sprintf("~%s:so:", taskName))), nil)
	if it.Next() == true {
		t.Fatal("there still has single object in database")
	}
	it.Release()

	it = resultDB.NewIterator(util.BytesPrefix([]byte(fmt.Sprintf("~%s:po:", taskName))), nil)
	if it.Next() == true {
		t.Fatal("there still has partial object in database")
	}
	it.Release()
}
