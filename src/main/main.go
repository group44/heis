package main

import (
	"fmt"
	//"net"
	//"../types"
	//"../com"
	"../order"
	"../elevator"
	//"time"
)



func main() {

	localOrders := order.NewLocalTable()
	elevator.Init()
	for {
		elevator.ControlStateMachine()
	}
	//go com.Run()
	//go elevator.init()
	//order.Run()

	/*
	//testMap := com.NewPeerMap()
	
	fmt.Println("Test variables created successfully")

	//go com.UpdatePeerMap(testMap, CART_ID, peerch) 
	
	lt := types.NewLocalTable()
	gt := types.NewGlobalTable()
	fmt.Println(lt[0])
	fmt.Println(gt[0])
	
		
	//data.ClaimOrder(testOrder, &testTable)
	//fmt.Println(testTable)
	

	//lconn.Close()
	//bconn.Close()
	*/

	fmt.Println("End")
}