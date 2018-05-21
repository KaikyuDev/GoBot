package gobot

import (
	"io/ioutil"
	"log"
)

func ReadFileBytes(fileName string) []byte {
	b, err := ioutil.ReadFile(fileName) // just pass the file name
	if err != nil {
		log.Println(err)
	}
	return b
}

func ReadFileString(fileName string) string {
	return string(ReadFileBytes(fileName))
}
