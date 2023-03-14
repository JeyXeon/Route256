package schema

type Reservation struct {
	OrderId     int64  `db:"order_id"`
	Sku         uint32 `db:"sku"`
	WarehouseId int64  `db:"warehouse_id"`
	Count       uint32 `db:"count"`
}
