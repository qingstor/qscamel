package generator

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func generate(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	ur, err := os.Open("/dev/urandom")
	if err != nil {
		log.Fatalf("Open urandom failed for %v.", err)
	}
	defer ur.Close()
	content := make([]byte, 512000)
	_, err = ur.Read(content)
	if err != nil {
		log.Fatalf("Read failed for %v.", err)
	}

	for j := 0; j < 1000; j++ {
		os.MkdirAll(fmt.Sprintf("/data/%d/%d", id, j), 0600)
		for k := 0; k < 15; k++ {
			f, err := os.OpenFile(
				fmt.Sprintf("/data/%d/%d/%d", id, j, k),
				os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				log.Fatalf("Open failed for %v.", err)
			}
			_, err = f.Write(content)
			if err != nil {
				log.Fatalf("Write failed for %v.", err)
			}

			f.Close()

			log.Printf("/data/%d/%d/%d created.", id, j, k)
		}
	}
}

func main() {
	wg := sync.WaitGroup{}

	for i := 0; i < 256; i++ {
		wg.Add(1)
		go generate(i, &wg)
	}

	wg.Wait()
}
