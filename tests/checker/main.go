package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	vis := [256][1000][15]int{}

	re, err := regexp.Compile("Object /(\\d+)/(\\d+)/(\\d+) is created.")
	if err != nil {
		log.Fatalf("Regexp compiled failed for %v.", err)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Open file %s failed for %v.", os.Args[1], err)
	}
	defer f.Close()

	buf := bufio.NewScanner(f)
	for buf.Scan() {
		t := buf.Bytes()

		if !re.Match(t) {
			continue
		}

		sub := re.FindAllSubmatch(t, -1)
		x, _ := strconv.Atoi(string(sub[0][1][:]))
		y, _ := strconv.Atoi(string(sub[0][2][:]))
		z, _ := strconv.Atoi(string(sub[0][3][:]))
		//log.Printf("x=%d, y=%d, z=%d.", x, y, z)
		if x > 255 {
			continue
		}
		if y > 999 {
			continue
		}
		if z > 14 {
			continue
		}
		vis[x][y][z] = 1
	}

	for i := 0; i < 256; i++ {
		for j := 0; j < 1000; j++ {
			for k := 0; k < 15; k++ {
				if vis[i][j][k] == 1 {
					continue
				}
				log.Printf("Object /%d/%d/%d is missing.", i, j, k)
			}
		}
	}
}
