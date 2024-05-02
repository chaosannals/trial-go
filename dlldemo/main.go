package main

import "C"

//export SumDemo
func SumDemo(a int, b int) int {
	return a + b
}

func main() {

}
