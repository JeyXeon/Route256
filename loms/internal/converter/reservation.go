package converter

import (
	"route256/loms/internal/model"
	"route256/loms/internal/repository/schema"
)

func ToReservationListModel(reservations []*schema.Reservation) []*model.Reservation {
	if reservations == nil {
		return nil
	}

	result := make([]*model.Reservation, 0, len(reservations))
	for _, reservation := range reservations {
		result = append(result, ToReservationModel(reservation))
	}
	return result
}

func ToReservationModel(reservation *schema.Reservation) *model.Reservation {
	if reservation == nil {
		return nil
	}

	return &model.Reservation{
		OrderId:     reservation.OrderId,
		Sku:         reservation.Sku,
		WareHouseId: reservation.WarehouseId,
		Count:       reservation.Count,
	}
}
