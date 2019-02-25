package utils

import (
	"log"
	"reflect"

	"github.com/yunify/qscamel/utils"
)

// CompareLocalDirectoryMD5 will compare all the md5
// of file in the directory, fatal if not equal
func CompareLocalDirectoryMD5(d1, d2 string) bool {
	kv1, err := utils.GetDirKvPair(d1)
	if err != nil {
		log.Fatal(err)
	}
	kv2, err := utils.GetDirKvPair(d2)
	if err != nil {
		log.Fatal(err)
	}

	return reflect.DeepEqual(kv1, kv2)
}

// CheckDirectroyEqual check two dirctory if is equal
func CheckDirectroyEqual(fmap map[string]string) {
	eq := CompareLocalDirectoryMD5(fmap["dir"]+"/src", fmap["dir"]+"/dst")
	if !eq {
		log.Fatal("check directory is finished: not equal")
	} else {
		log.Println("check directory is finished: equal")
	}


}
