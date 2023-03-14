package schema

type OrderStatus string

const (
	New             OrderStatus = "NEW"
	AwaitingPayment OrderStatus = "AWAITING_PAYMENT"
	Failed          OrderStatus = "FAILED"
	Payed           OrderStatus = "PAYED"
	Cancelled       OrderStatus = "CANCELLED"
)

type Order struct {
	Id     int64       `db:"order_id"`
	UserId int64       `db:"user_id"`
	Status OrderStatus `db:"status" sql:"type:order_status"`
}
