package model

type Reservation struct {
	OrderId     int64
	Sku         uint32
	WareHouseId int64
	Count       uint32
}
