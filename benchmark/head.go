package main

import (
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	var fl *os.File
	var err error
	if len(os.Args) >= 3 {
		name := os.Args[2]
		fl, err = os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fl = os.Stdin
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1024*1024)
	for n > len(buf) {
		k, err := io.ReadFull(fl, buf)
		if err != nil {
			return
		}
		n -= k
	}
	io.ReadFull(fl, buf[:n])
}
