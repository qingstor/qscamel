package utils

import (
	"crypto/md5"
	"io/ioutil"
	"os"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/qiniu/x/log.v7"
)

const (
	dirworker  = 2
	crypworker = 3
)

// GetDirKvPair will iterates the whole `dir`, return a kv map
// the key is position in the directory, and the value is the
// md5 sum of the file content.
func GetDirKvPair(dir string) (map[string]string, error) {
	wg := sync.WaitGroup{}
	dirch := make(chan []string, 1)
	var listdone int32
	fch := make(chan string)
	errch := make(chan error, 1)

	md5Pair := make(map[string]string)
	var Lock sync.Mutex

	dirch <- []string{dir}

	// list dir
	for i := 0; i < dirworker; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case dirsli := <-dirch:
					var dirlist []string
					for _, dir := range dirsli {
						dp, err := ioutil.ReadDir(dir)
						if err != nil {
							errch <- err
							break
						}

						for _, et := range dp {
							if et.IsDir() {
								dirlist = append(dirlist, dir+"/"+et.Name())
								continue
							}
							fch <- dir + "/" + et.Name()
						}
					}
					if len(dirlist) != 0 {
						dirch <- dirlist
					}
				default:
					atomic.AddInt32(&listdone, 1)
					wg.Done()
					return
				}
			}
		}()
	}

	// cryption
	for i := 0; i < crypworker; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case f := <-fch:
					fp, err := os.Open(f)
					if err != nil {
						log.Info(err)
						errch <- err
						continue
					}
					cnt, err := ioutil.ReadAll(fp)
					fp.Close()
					if err != nil {
						log.Info(err)
					}
					h := md5.New()
					h.Write(cnt)
					name := f[len(dir)+1:]
					Lock.Lock()
					md5Pair[name] = string(h.Sum(nil))
					Lock.Unlock()
				default:
					if listdone == dirworker {
						wg.Done()
						return
					}

				}
			}
		}()
	}

	wg.Wait()
	select {
	case err := <-errch:
		return nil, err
	default:
	}
	close(fch)
	close(dirch)
	return md5Pair, nil
}

// CompareLocalDirectoryMD5 will compare all the md5
// of file in the directory, return true is equal
func CompareLocalDirectoryMD5(d1, d2 string) (bool, error) {
	kv1, err := GetDirKvPair(d1)
	if err != nil {
		return false, err
	}
	kv2, err := GetDirKvPair(d2)
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(kv1, kv2), nil
}
