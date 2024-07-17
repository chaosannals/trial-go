package main

import (
	"encoding/json"
	"log"
	"time"
)

type TS[T any] struct {
	Id     string    `json:"id"`
	Token  T         `json:"token"`
	TimeAt time.Time `json:"timeAt"`
}

func gT[T any](t T) {
	ts := TS[T]{
		Id:     "aaa",
		Token:  t,
		TimeAt: time.Now(),
	}
	log.Printf("ts: %v\n", ts)
	tb, err := json.Marshal(&ts)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("tb: %v\n", tb)
	tbs := string(tb)
	log.Printf("tbs: %s\n", tbs)

	ts2 := TS[T]{}
	if err := json.Unmarshal(tb, &ts2); err != nil {
		log.Fatalln(err)
	}
	log.Printf("ts2: %v", ts2)
}

func main() {
	gT("BBB")
}
