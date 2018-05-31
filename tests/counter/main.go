package main

import (
	"io/ioutil"
	"log"
	"os"
)

func count(f os.FileInfo, cur string) int64 {
	if !f.IsDir() {
		return 1
	}

	total := int64(0)
	fs, err := ioutil.ReadDir(cur + f.Name())
	if err != nil {
		log.Fatalf("Readdir failed for %v.", err)
	}

	for _, v := range fs {
		total += count(v, cur+f.Name()+"/")
	}

	log.Printf("Counted %d.", total)
	return total
}

func main() {
	f, err := os.Stat(os.Args[1])
	if err != nil {
		log.Fatalf("Sata file %s failed for %v.", os.Args[1], err)
	}

	log.Printf("Total files are %d.", count(f, "/"))
}
