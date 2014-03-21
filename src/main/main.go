package main

import (
	"fmt"
	//"net"
	//"../types"
	"../com"
	"../order"
	//"../driver"
	//"../elevator"
	//"time"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	done := make(chan bool)

	go order.Run()
	go com.Run()
	//go elevator.Run()
	

	<-done
	fmt.Println("End")
}
