package data

const CART_ID int = 1

type Order struct {
	Floor, Dir, Cart int
}

func NewOrder(f int, d int, c int) Order {
	return Order{f, d, c}
}

type OrderTable [4][2]Order

func ClaimOrder(o Order, t *OrderTable) {
	t[o.Floor][o.Dir].Cart = CART_ID
}

func InsertOrder(o Order, t *OrderTable) {
	t[o.Floor][o.Dir] = o
}


