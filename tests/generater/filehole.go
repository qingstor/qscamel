package generater

import (
	"io/ioutil"
	"os"
)

// CreateHoleFile create numbers of file with hole in pointed directory
func CreateHoleFile(dir string, fileSize, holeSize, offset int64, n int) error {
	for i := 0; i < n; i++ {
		file, err := ioutil.TempFile(dir+"/", "FILE*.hole")
		if err != nil {
			return err
		}
		content, err := CreateRandomByteStream(offset)
		if err != nil {
			return err
		}

		if _, err = file.Write(content); err != nil {
			return err
		}

		if _, err = file.Seek(holeSize, os.SEEK_SET); err != nil {
			return err
		}

		content, err = CreateRandomByteStream(fileSize - offset)
		if err != nil {
			return err
		}

		if _, err = file.Write(content); err != nil {
			return err
		}
		file.Close()
	}
	return nil
}
