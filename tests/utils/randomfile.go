package utils

import (
	"crypto/rand"
	"io/ioutil"
	"log"
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

// piece used to allocate bytes slice
const piece = 256 * MB

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
func CreateRandomByteStream(stm *[]byte) error {
	_, err := rand.Read(*stm)
	if err != nil {
		return err
	}
	return nil
}

// CreateTestRandomFile create the `filePerDir` number of random file
// in the `dir` directory.
func CreateTestRandomFile(filePerDir int, fileSize int64, dir string) error {
	for i := 0; i < filePerDir; i++ {
		file, err := ioutil.TempFile(dir+"/", "TESTFILE*.camel")
		if err == nil {
			content := make([]byte, piece)
			bs := fileSize
			// avoid the slice allocation out of memory
			for bs = fileSize; bs > piece; bs -= piece {
				err := CreateRandomByteStream(&content)
				if err == nil {
					_, err := file.Write(content)
					err = file.Sync()
					if err == nil {
						continue
					}
				}
				return err
			}

			content = make([]byte, bs)
			err := CreateRandomByteStream(&content)
			if err == nil {
				_, err := file.Write(content)
				if err == nil {
					file.Close()
					continue
				}
			}
		}
		return err
	}
	return nil
}

// CreateHoleFile create numbers of file with hole in pointed directory
func CreateHoleFile(dir string, fileSize, holeSize, offset int64, n int) {
	for i := 0; i < n; i++ {
		file, err := ioutil.TempFile(dir+"/", "FILE*.hole")
		if err != nil {
			log.Fatal(err)
		}
		content := make([]byte, offset)
		err = CreateRandomByteStream(&content)
		if err != nil {
			log.Fatal(err)
		}

		if _, err = file.Write(content); err != nil {
			log.Fatal(err)
		}

		if _, err = file.Seek(holeSize, 0); err != nil {
			log.Fatal(err)
		}

		content = make([]byte, fileSize-offset)
		err = CreateRandomByteStream(&content)
		if err != nil {
			log.Fatal(err)
		}

		if _, err = file.Write(content); err != nil {
			log.Fatal(err)
		}
		file.Close()
	}
}
