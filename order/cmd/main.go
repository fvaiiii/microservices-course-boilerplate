package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/fvaiiii/microservices-course-boilerplate/order/pkg/app"
	inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/payment/v1"
)

const (
	httpAddress             = ":8080"
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"
	grpcKeepaliveTime       = 5 * time.Minute
	grpcKeepaliveTimeout    = 1 * time.Second
	httpReadHeaderTimeout   = 5 * time.Second
	httpReadTimeout         = 15 * time.Second
	httpWriteTimeout        = 15 * time.Second
	httpIdleTimeout         = 60 * time.Second
	httpShutdownTimeout     = 30 * time.Second
	httpMaxHeaderBytes      = 1 << 20
)

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepaliveTime,
			Timeout:             grpcKeepaliveTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		os.Exit(1)
	}

	paymentConn, err := grpc.NewClient(
		paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepaliveTime,
			Timeout:             grpcKeepaliveTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		if closeErr := inventoryConn.Close(); closeErr != nil {
			slog.Error("ошибка закрытия inventoryConn", "error", closeErr)
		}
		os.Exit(1)
	}
	inventoryClient := inventoryv1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentv1.NewPaymentServiceClient(paymentConn)

	handler, err := app.NewHTTPHandler(inventoryClient, paymentClient)
	if err != nil {
		slog.Error("ошибка создания хендлера", "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:              httpAddress,
		Handler:           handler,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
		MaxHeaderBytes:    httpMaxHeaderBytes,
	}

	go func() {
		slog.Info("запуск OrderService", "port", 8080)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("ошибка запуска сервера", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("ошибка graceful shutdown", "error", err)
	}

	if err := inventoryConn.Close(); err != nil {
		slog.Error("ошибка закрытия inventoryConn", "error", err)
	}

	if err := paymentConn.Close(); err != nil {
		slog.Error("ошибка закрытия paymentConn", "error", err)
	}

	slog.Info("OrderService остановлен")
}
