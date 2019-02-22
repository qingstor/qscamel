package utils

import (
	"reflect"
	"testing"

	"github.com/yunify/qscamel/utils"
)

// CompareLocalDirectoryMD5 will compare all the md5
// of file in the directory, fatal if not equal
func CompareLocalDirectoryMD5(t testing.TB, d1, d2 string) bool {
	kv1, err := utils.GetDirKvPair(d1)
	if err != nil {
		t.Fatal(err)
	}
	kv2, err := utils.GetDirKvPair(d2)
	if err != nil {
		t.Fatal(err)
	}

	return reflect.DeepEqual(kv1, kv2)
}

// CheckDirectroyEqual check two dirctory if is equal
func CheckDirectroyEqual(t testing.TB, fmap map[string]string) {
	eq := CompareLocalDirectoryMD5(t, fmap["dir"]+"/src", fmap["dir"]+"/dst")
	if !eq {
		t.Error("check directory is finished: not equal")
	} else {
		t.Log("check directory is finished: equal")
	}


}
