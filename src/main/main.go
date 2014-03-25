package main

import (
	"../com"
	"../elevator"
	"../order"
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	done := make(chan bool)

	go order.Run()
	go com.Run()
	go elevator.Run()
	go order.OsTest()

	<-done
	fmt.Println("End")
}
