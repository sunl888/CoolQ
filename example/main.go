package main

import (
	"fmt"
	"time"
)

func main() {
	fun1()
	fmt.Println("hello world")
	time.Sleep(time.Second*8)
}

func fun1(){
	go func() {
		time.Sleep(time.Second*4)
		fmt.Println("hello Goroutine")
	}()
}
