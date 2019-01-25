package generater

import (
	"crypto/rand"
	"io/ioutil"
	"math"
)

// B means Bytes
const B int64 = 1

// KB means kilobytes
const KB = B * 1024

// MB means megabytes
const MB = KB * 1024

// GB means gigabytes
const GB = MB * 1024

// caculate Geometric series sum
func seriesSum(dirnum, depth int) int {
	if dirnum == 0 {
		return 0
	} else if dirnum == 1 {
		return depth * dirnum
	}
	return dirnum * int(((1 - math.Pow(float64(dirnum), float64(depth))) / float64(1-dirnum)))
}

// CreateRandomByteStream return `size` of random bytes
func CreateRandomByteStream(size int64) ([]byte, error) {
	p := make([]byte, size)
	_, err := rand.Read(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// CreateTestRandomFile create the `filePerDir` number of random file
// in the `dir` directory.
func CreateTestRandomFile(filePerDir int, fileSize int64, dir string) error {
	for i := 0; i < filePerDir; i++ {
		file, err := ioutil.TempFile(dir+"/", "TESTFILE*.camel")
		if err == nil {
			content, err := CreateRandomByteStream(fileSize)
			if err == nil {
				_, err := file.Write(content)
				if err == nil {
					defer file.Close()
					continue
				}
			}
		}
		return err
	}
	return nil
}
