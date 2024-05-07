package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func makeSuffleRange(min int, max int) []int {
	c := max - min
	r := make([]int, c)
	for i := 0; i < c; i++ {
		r[i] = min + i
	}
	rand.Shuffle(c, func(a int, b int) {
		r[a], r[b] = r[b], r[a]
	})
	return r
}

// 多个写入会导致 fatal error: concurrent map writes
func testLangMap() {
	m := make(map[string]int)
	c := make(chan int, 1)
	for j := 0; j <= 10; j++ {
		k := j
		go func() {
			min := rand.Intn(100)
			vs := makeSuffleRange(min, min+1000)
			fmt.Printf("a: %d %v\n", k, vs)
			for i := range vs {
				m["a"] = j*100 + i
				time.Sleep(time.Millisecond)
			}
		}()
		go func() {
			min := rand.Intn(100)
			vs := makeSuffleRange(min, min+1000)
			fmt.Printf("b: %d %v\n", k, vs)
			for i := range vs {
				m["a"] = j*100 + i
				time.Sleep(time.Millisecond)
			}
		}()
	}
	go func() {
		for i := 0; i <= 1000; i++ {
			fmt.Println(m["a"])
			time.Sleep(time.Millisecond)
		}
		c <- -1
	}()

	r := <-c
	fmt.Println(r)
}

// 
func testSyncMap() {
	m := sync.Map{}
	c := make(chan int, 1)
	for j := 0; j <= 10; j++ {
		k := j
		go func() {
			min := rand.Intn(100)
			vs := makeSuffleRange(min, min+1000)
			fmt.Printf("a: %d %v\n", k, vs)
			for i := range vs {
				m.Store("a", j*100+i)
				time.Sleep(time.Millisecond)
			}
		}()
		go func() {
			min := rand.Intn(100)
			vs := makeSuffleRange(min, min+1000)
			fmt.Printf("b: %d %v\n", k, vs)
			for i := range vs {
				m.Store("a", j*100+i)
				time.Sleep(time.Millisecond)
			}
		}()
	}
	go func() {
		for i := 0; i <= 1000; i++ {
			if v, ok := m.Load("a"); ok {
				fmt.Println(v)
			} else {
				fmt.Println("Load a failed")
			}
			time.Sleep(time.Millisecond)
		}
		c <- -1
	}()

	r := <-c
	fmt.Println(r)
}

func main() {
	// testLangMap()
	testSyncMap()
}
