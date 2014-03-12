package main

import (
	"fmt"
	"./driver"
	"./order"
	"./elevator"
)



func main() {
	
	order.PrintTable()
	fmt.Println(driver.ElevGetButtonSignal(2, 2))
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.UpdateLocalTable()
	order.PrintTable()
	driver.ElevInit()
	order.Init()
	elevator.Init()// litt dumt med likt nanv?
	
	for{
		order.UpdateLocalTable()
		driver.ElevSetFloorIndicator(driver.ElevGetFloorSensorSignal())
		elevator.ControlStateMachine()
	}
	
	
}
