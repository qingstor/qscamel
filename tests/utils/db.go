package utils

import (
	"fmt"
	"testing"

	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/yunify/qscamel/db"
)

// CheckDBEmpty check the temp dbfile and testing
// whether the Database is empty.
func CheckDBEmpty(t *testing.T, fmap *map[string]string) {
	db, err := db.NewDB(&db.DatabaseOptions{(*fmap)["dir"] + "/db"})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	it := db.NewIterator(nil, nil)
	if it.Next() == true {
		t.Fatal("database is not empty")
	}
}

// CheckDBNoObject check the temp dbfile and testing
// whether the Database has any object
func CheckDBNoObject(t *testing.T, fmap *map[string]string) {
	db, err := db.NewDB(&db.DatabaseOptions{(*fmap)["dir"] + "/db"})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	taskName := (*fmap)["name"]
	it := db.NewIterator(util.BytesPrefix([]byte(fmt.Sprintf("~%s:do:", taskName))), nil)
	if it.Next() == true {
		t.Fatal("there still has directory object in database")
	}
	it.Release()

	it = db.NewIterator(util.BytesPrefix([]byte(fmt.Sprintf("~%s:so:", taskName))), nil)
	if it.Next() == true {
		t.Fatal("there still has single object in database")
	}
	it.Release()

	it = db.NewIterator(util.BytesPrefix([]byte(fmt.Sprintf("~%s:po:", taskName))), nil)
	if it.Next() == true {
		t.Fatal("there still has partial object in database")
	}
	it.Release()
}
