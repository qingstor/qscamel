package utils

import (
	"testing"

	"github.com/magiconair/properties/assert"

	"github.com/yunify/qscamel/utils"
)

// CompareLocalDirectoryMD5 will compare all the md5
// of file in the directory, fatal if not equal
func CompareLocalDirectoryMD5(t *testing.T, d1, d2 string) {
	kv1, err := utils.GetDirKvPair(d1)
	if err != nil {
		t.Fatal(err)
	}
	kv2, err := utils.GetDirKvPair(d2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, kv1, kv2, "check directory is finished: not equal")

}

// CheckDirectroyEqual check two dirctory if is equal
func CheckDirectroyEqual(t *testing.T, fmap *map[string]string) {
	CompareLocalDirectoryMD5(t, (*fmap)["dir"]+"/src", (*fmap)["dir"]+"/dst")
	t.Logf("check directory is finished: equal")
}
