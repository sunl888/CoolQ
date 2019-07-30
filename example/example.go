package main

import "fmt"
import "time"

func main1() {

	// make the request chan chan that both go-routines will be given
	requestChan := make(chan chan string)

	// start the goroutines
	go goroutineC(requestChan)
	go goroutineD(requestChan)

	// sleep for a second to let the goroutines complete
	time.Sleep(time.Second * 4)

}

func goroutineC(requestChan chan chan string) {

	// make a new response chan
	responseChan := make(chan string)

	// send the responseChan to goRoutineD
	requestChan <- responseChan

	// read the response
	response := <-responseChan

	fmt.Printf("Response: %v\n", response)

}

func goroutineD(requestChan chan chan string) {

	// read the responseChan from the requestChan
	responseChan := <-requestChan
	time.Sleep(time.Second * 2)
	// send a value down the responseChan
	responseChan <- "wassup!"

}
