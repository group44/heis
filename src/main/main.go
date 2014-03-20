package main

import (
	"fmt"
	//"net"
	//"../types"
	"../com"
	"../order"
	//"../driver"
	"../elevator"
	//"time"
	"runtime"
)

const CART_ID int = 1

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	done := make(chan bool)

	go order.Run()
	go com.Run()
	go elevator.Run()

	/*
	   //Todo
	   go order.Run()
	*/

	<-done
	fmt.Println("End")
}
