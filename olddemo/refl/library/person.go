package library

import (
	"fmt"
)

type Person struct {
	name string
}

func (i *Person) SayByPointer(msg string) {
	fmt.Println(msg)
}

func (i *Person) WalkByPointer() {

}

func (i Person) SayByValue(msg string) {
	fmt.Println(msg)
}

func (i Person) WalkByValue() {

}
