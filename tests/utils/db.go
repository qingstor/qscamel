package utils

import (
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/yunify/qscamel/db"
)

// CheckDBEmpty check the temp dbfile and testing
// whether the Database is empty.
func CheckDBEmpty(fmap map[string]string) {
	resultDB, err := db.NewDB(&db.DatabaseOptions{Address: fmap["dir"] + "/db"})
	if err != nil {
		log.Fatal(err)
	}
	defer resultDB.Close()
	it := resultDB.NewIterator(nil, nil)
	if it.Next() == true {
		log.Fatal("database is not empty")
	}
	it.Release()
}

// CheckDBNoObject check the temp dbfile and testing
// whether the Database has any object
func CheckDBNoObject(fmap map[string]string) {
	resultDB, err := db.NewDB(&db.DatabaseOptions{Address: fmap["dir"] + "/db"})
	if err != nil {
		log.Fatal(err)
	}
	defer resultDB.Close()

	taskName := fmap["name"]
	it := resultDB.NewIterator(util.BytesPrefix([]byte(fmt.Sprintf("~%s:do:", taskName))), nil)
	if it.Next() == true {
		log.Fatal("there still has directory object in database")
	}
	it.Release()

	it = resultDB.NewIterator(util.BytesPrefix([]byte(fmt.Sprintf("~%s:so:", taskName))), nil)
	if it.Next() == true {
		log.Fatal("there still has single object in database")
	}
	it.Release()

	it = resultDB.NewIterator(util.BytesPrefix([]byte(fmt.Sprintf("~%s:po:", taskName))), nil)
	if it.Next() == true {
		log.Fatal("there still has partial object in database")
	}
	it.Release()
}
