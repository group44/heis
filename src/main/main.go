package main

import (
	"fmt"
	//"net"
	//"../types"
	"../com"
	//"../order"
	//"../driver"
	//"../elevator"
	//"time"
)

const CART_ID int = 1

func main() {

    done := make(chan bool)

    //go order.Run()
    go com.Run()
    //go elevator.Run()

    /*
    //Todo
    go order.Run()
    */

    <- done
    fmt.Println("End")
}