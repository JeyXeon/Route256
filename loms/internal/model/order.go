package model

type OrderStatus string

const (
	New             OrderStatus = "new"
	AwaitingPayment OrderStatus = "awaiting_payment"
	Failed          OrderStatus = "failed"
	Payed           OrderStatus = "payed"
	Cancelled       OrderStatus = "cancelled"
)

type OrderItem struct {
	Sku   uint32
	Count uint32
}

type Order struct {
	ID     int64
	User   int64
	Status OrderStatus
	Items  []*OrderItem
}
