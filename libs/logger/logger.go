package logger

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var globalLogger *zap.Logger

func Init(devel bool) {
	globalLogger = New(devel)
}

func New(devel bool) *zap.Logger {
	var logger *zap.Logger
	var err error
	if devel {
		logger, err = zap.NewDevelopment()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger, err = cfg.Build()
	}
	if err != nil {
		panic(err)
	}

	return logger
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := handler(ctx, req)

	if err != nil {
		globalLogger.Warn(
			"gRPC",
			zap.String("server", info.FullMethod),
			zap.String("method", info.FullMethod),
			zap.String("err", err.Error()),
		)
		return nil, err
	}

	globalLogger.Debug(
		"gRPC",
		zap.String("server", info.FullMethod),
		zap.String("method", info.FullMethod),
	)

	return res, nil
}

func GetLogger() *zap.Logger {
	return globalLogger
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}
