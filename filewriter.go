package main

import (
	"log"
	"os"
)

//File holds filename to open
type File struct {
	Filename string
}

//Write creates and/or appends to a file
func (f *File) Write(message []byte) {
	file, err := os.OpenFile(f.Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	_, err = file.Write(message)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}
