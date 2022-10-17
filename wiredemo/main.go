package main

import (
	"log"
)


func main() {
	ts, err := InitTrial()
	if err != nil {
		log.Fatal(err)
	}
	if err := ts.Serve(); err != nil {
		log.Fatal(err)
	}
}