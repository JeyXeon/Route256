package schema

type Stock struct {
	Sku         uint32 `db:"sku"`
	WarehouseId int64  `db:"warehouse_id"`
	Count       uint64 `db:"count"`
}
