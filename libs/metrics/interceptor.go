package metrics

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

func ServerMetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ServerRequestsCounter.Inc()

	timeStart := time.Now()

	res, err := handler(ctx, req)

	elapsed := time.Since(timeStart)

	resStatus, _ := status.FromError(err)

	ServerHistogramResponseTime.WithLabelValues(fmt.Sprint(resStatus.Code())).Observe(elapsed.Seconds())
	ServerResponseCounter.WithLabelValues(fmt.Sprint(resStatus.Code())).Inc()

	return res, err
}

func ClientMetricsInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	timeStart := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)

	elapsed := time.Since(timeStart)

	resStatus, _ := status.FromError(err)

	ClientHistogramResponseTime.WithLabelValues(fmt.Sprint(resStatus.Code())).Observe(elapsed.Seconds())

	return err
}
