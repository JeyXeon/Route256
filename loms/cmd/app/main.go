package main

import (
	"context"
	"net"
	"route256/libs/dbmanager"
	"route256/libs/kafka"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/libs/tracing"
	loms "route256/loms/internal/api"
	"route256/loms/internal/config"
	"route256/loms/internal/repository/postgres"
	"route256/loms/internal/service"
	desc "route256/loms/pkg/loms"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger.Init(false)

	err := config.Init()
	if err != nil {
		logger.Fatal("config init", zap.Error(err))
	}

	port := config.ConfigData.Port

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	tracing.Init(logger.GetLogger(), "loms")
	metrics.Init()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
				metrics.ServerMetricsInterceptor,
				logger.LoggingInterceptor,
			),
		),
	)
	reflection.Register(s)

	grpc_prometheus.Register(s)
	go metrics.ServeMetrics(config.ConfigData.MetricsPort, logger.GetLogger())

	lomsDbUrl := config.ConfigData.LomsDbUrl
	pool, err := pgxpool.Connect(context.Background(), lomsDbUrl)
	if err != nil {
		logger.Fatal("db connect", zap.Error(err))
	}

	dbManager := dbmanager.New(pool)

	brokers := config.ConfigData.Kafka.Brokers
	kafkaSyncProducer, err := kafka.NewSyncProducer(brokers)
	if err != nil {
		logger.Fatal("kafka sync producer", zap.Error(err))
	}

	reservationsRepository := postgres.NewReservationsRepository(dbManager)
	stocksRepository := postgres.NewStocksRepository(dbManager)
	orderRepository := postgres.NewOrderRepository(dbManager)
	outboxKafkaRepository := postgres.NewKafkaOutboxRepository(dbManager)

	lomsService := service.New(
		dbManager,
		kafkaSyncProducer,
		reservationsRepository,
		stocksRepository,
		orderRepository,
		outboxKafkaRepository,
	)

	go lomsService.CheckPaymentTimeoutCron(context.Background())
	go lomsService.RunRecordProcessor(context.Background())
	go lomsService.RunRecordCleaner(context.Background())

	desc.RegisterLomsServer(s, loms.NewLoms(lomsService))

	logger.Info("server listening", zap.Any("address", lis.Addr()))

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
