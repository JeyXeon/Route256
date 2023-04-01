package service

import (
	"context"
	"route256/loms/internal/model"
	"time"

	"github.com/Shopify/sarama"
)

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
}

type KafkaSyncProducer interface {
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type ReservationsRepository interface {
	GetReservations(ctx context.Context, orderId int64) ([]*model.Reservation, error)
	AddReservations(ctx context.Context, orderItems []*model.Reservation) error
	RemoveReservations(ctx context.Context, orderId int64) error
	RemoveReservationsByOrderIds(ctx context.Context, orderId []int64) error
}

type StocksRepository interface {
	GetStocks(ctx context.Context, skus []uint32) ([]*model.Stock, error)
	RevertReservations(ctx context.Context, reservations []*model.Reservation) error
	WriteOffStocks(ctx context.Context, stocks []*model.Stock) error
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, userId int64) (int64, error)
	GetOrder(ctx context.Context, orderId int64) (*model.Order, error)
	GetTimeoutedPaymentOrderIds(ctx context.Context, time time.Time) ([]int64, error)
	UpdateOrderStatus(ctx context.Context, orderId int64, newStatus model.OrderStatus) error
	UpdateOrdersStatuses(ctx context.Context, orderIds []int64, newStatus model.OrderStatus) ([]int64, error)
}

type OutboxKafkaRepository interface {
	GetUnprocessedRecords(ctx context.Context) ([]*model.KafkaRecord, error)
	UpdateRecordByID(ctx context.Context, message *model.KafkaRecord) error
	RemoveRecordsBeforeDatetime(ctx context.Context, expireTime time.Time) error
	CreateKafkaRecord(ctx context.Context, record *model.KafkaRecord) error
}

type Service struct {
	transactionManager     TransactionManager
	kafkaSyncProducer      KafkaSyncProducer
	reservationsRepository ReservationsRepository
	stocksRepository       StocksRepository
	orderRepository        OrderRepository
	outboxKafkaRepository  OutboxKafkaRepository
}

func New(
	transactionManager TransactionManager,
	kafkaSyncProducer KafkaSyncProducer,
	reservationsRepository ReservationsRepository,
	stocksRepository StocksRepository,
	orderRepository OrderRepository,
	outboxKafkaRepository OutboxKafkaRepository,
) *Service {
	return &Service{
		transactionManager:     transactionManager,
		kafkaSyncProducer:      kafkaSyncProducer,
		reservationsRepository: reservationsRepository,
		stocksRepository:       stocksRepository,
		orderRepository:        orderRepository,
		outboxKafkaRepository:  outboxKafkaRepository,
	}
}
