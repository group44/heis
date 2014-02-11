package elev



type order struct {
	floor int
	dir int
}

func NewOrder(f int, d int) order {
	return order{f, d}
}
