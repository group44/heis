package elev



type Order struct {
	floor int
	dir int
}

func NewOrder(f int, d int) Order {
	return Order{f, d}
}
