package order

/*
import (
	"../types"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)


var (
	outputCh            = make(chan types.Data, 5)
	addOrderCh          = make(chan []int, 5)
	removeOrderCh       = make(chan []int, 5)
	updateGlobalTableCh = make(chan types.GlobalTable) //få buffer på denne
	done                = make(chan bool)
	GlobalOrders        types.GlobalTable
)

func main() {
	GlobalOrders = NewGlobalTable()
	brdAddr := "localhost:12345"
	lisAddr := ":12346"

	lAddr, err := net.ResolveUDPAddr("udp", lisAddr)
	CheckError(err)
	bAddr, err := net.ResolveUDPAddr("udp", brdAddr)
	CheckError(err)

	bConn, err := net.DialUDP("udp", nil, bAddr)
	CheckError(err)
	lConn, err := net.ListenUDP("udp", lAddr)
	CheckError(err)
	time.Sleep(5000 * time.Millisecond)
	fmt.Println("hei verden")

	go Run()
	go Recive(lConn)
	go Cast(bConn)
	go AddOrder()
	go RemoveOrder()

	<-done
}

func NewGlobalTable() [][]int {
	t := make([][]int, types.N_FLOORS)
	for i := range t {
		t[i] = make([]int, 2)
		for j := range t[i] {
			t[i][j] = 0
		}
	}
	return t
}

func Run() {

	var cost int
	order := make([]int, 2)
	time.Sleep(25 * time.Millisecond)

	for {
		fmt.Println("test")
		order[0] = 9
		order[1] = 9
		outputCh <- types.Data{Head: "order", Order: order, Cost: cost, Table: GlobalOrders}
		time.Sleep(2000 * time.Millisecond)

		cost = 55
		outputCh <- types.Data{Head: "cost", Order: order, Cost: cost, Table: GlobalOrders}
		time.Sleep(2000 * time.Millisecond)

		outputCh <- types.Data{Head: "table", Order: order, Cost: cost, Table: GlobalOrders}
		time.Sleep(2000 * time.Millisecond)
		order[0] = 1
		order[1] = 1

		fmt.Println("Global table is:")
		fmt.Println(GlobalOrders)
		outputCh <- types.Data{Head: "addorder", Order: order, Cost: cost, Table: GlobalOrders}

		time.Sleep(1 * time.Millisecond)

		order[0] = 3
		order[1] = 1

		fmt.Println("Global table is:")
		fmt.Println(GlobalOrders)
		outputCh <- types.Data{Head: "addorder", Order: order, Cost: cost, Table: GlobalOrders}
		time.Sleep(2000 * time.Millisecond)
		fmt.Println("Global table is:")
		fmt.Println(GlobalOrders)

		outputCh <- types.Data{Head: "removeorder", Order: order, Cost: cost, Table: GlobalOrders}
		time.Sleep(2000 * time.Millisecond)

	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func Cast(outConn *net.UDPConn) {
	fmt.Println("Cast")
	var out types.Data
	var buf = make([]byte, 1024)
	var err error

	for {
		time.Sleep(25 * time.Millisecond)
		out = <-outputCh
		out.ID = types.CART_ID
		out.T = time.Now()

		fmt.Println("Data casted:")
		fmt.Println(out)

		buf, err = json.Marshal(out)
		CheckError(err)
		_, err = outConn.Write(buf)
		CheckError(err)
	}
}

func Recive(inConn *net.UDPConn) {
	fmt.Println("Recieve")
	var inc types.Data

	var buf = make([]byte, 1024)

	for {
		time.Sleep(25 * time.Millisecond)
		n, _, err := inConn.ReadFromUDP(buf[0:])
		CheckError(err)
		err = json.Unmarshal(buf[:n], &inc)
		CheckError(err)

		switch inc.Head {

		case "order":
			//OrderCh <- inc.Order

			fmt.Println("Order received:")
			fmt.Println(inc.Order)
			fmt.Println("")

		case "table":
			updateGlobalTableCh <- inc.Table

			fmt.Println("Table received and updated")
			fmt.Println(inc.Table)
			fmt.Println("")

		case "cost":
			//fmt.Println(inc)
			//AuctionCh <- inc
			//fmt.Println(AuctionCh)

			fmt.Println("Cost received:")
			fmt.Println(inc.Cost)
			fmt.Println("")

		case "addorder":
			fmt.Println("Order added:")
			fmt.Println(inc.Order)
			fmt.Println("")
			addOrderCh <- inc.Order

		case "removeorder":
			removeOrderCh <- inc.Order
			fmt.Println("Order removed:")
			fmt.Println(inc.Order)
			fmt.Println("")

		default:

			fmt.Println("Default case entered")
		}
	}
}

//add order to the global table, must be fixed for the right cart number
//after order is added it sends the new table too be casted
func AddOrder() {
	order := make([]int, 2)
	for {
		order = <-addOrderCh
		cost := 0
		GlobalOrders[order[0]][order[1]] = 1
		outputCh <- types.Data{Head: "table", Order: order, Cost: cost, Table: GlobalOrders}
		fmt.Println("new global table:")
		fmt.Println(GlobalOrders)
		time.Sleep(25 * time.Millisecond)
	}
}

//removes order from globaltable, must be fixed for right cart number?
//after order is removed it sends the new table to be casted
func RemoveOrder() {
	order := make([]int, 2)
	for {
		order = <-removeOrderCh
		cost := 0
		GlobalOrders[order[0]][order[1]] = 0
		outputCh <- types.Data{Head: "table", Order: order, Cost: cost, Table: GlobalOrders}
		fmt.Println("new global table removed:")
		fmt.Println(GlobalOrders)
		time.Sleep(25 * time.Millisecond)
	}
}

//takes the casted table and updates the global table so all elevators have the same globaltable
func UpdateGlobalTable() {
	var tempTable = NewGlobalTable()
	for {
		tempTable = <-updateGlobalTableCh
		GlobalOrders = tempTable
		time.Sleep(25 * time.Millisecond)

	}
}
*/
