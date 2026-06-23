package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderrepo "github.com/fvaiiii/microservices-course-boilerplate/order/internal/repository/order"
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
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool, txManager, err := initDB(ctx)
	if err != nil {
		exitWithError(cancel, "инициализация БД", err)
	}
	defer pool.Close()

	inventoryConn, paymentConn, err := initGRPCClients()
	if err != nil {
		exitWithError(cancel, "инициализация gRPC-клиентов", err)
	}
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			slog.Error("ошибка закрытия inventoryConn", "error", err)
		}
		if err := paymentConn.Close(); err != nil {
			slog.Error("ошибка закрытия paymentConn", "error", err)
		}
	}()

	inventoryClient := inventoryv1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentv1.NewPaymentServiceClient(paymentConn)

	handler, err := app.NewHTTPHandler(pool, txManager, inventoryClient, paymentClient)
	if err != nil {
		exitWithError(cancel, "ошибка создания хендлера", err)
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

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("ошибка graceful shutdown", "error", err)
	}

	slog.Info("OrderService остановлен")
}

func initDB(ctx context.Context) (*pgxpool.Pool, orderrepo.TxManager, error) {
	pool, err := pgxpool.New(ctx, os.Getenv("DB_URI"))
	if err != nil {
		return nil, nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, err
	}

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		pool.Close()
		return nil, nil, err
	}

	slog.Info("подключение к PostgreSQL установлено")
	return pool, txManager, nil
}

func initGRPCClients() (*grpc.ClientConn, *grpc.ClientConn, error) {
	keepaliveParams := keepalive.ClientParameters{
		Time:                grpcKeepaliveTime,
		Timeout:             grpcKeepaliveTimeout,
		PermitWithoutStream: true,
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepaliveParams),
	}

	inventoryConn, err := grpc.NewClient(inventoryServiceAddress, opts...)
	if err != nil {
		return nil, nil, err
	}

	paymentConn, err := grpc.NewClient(paymentServiceAddress, opts...)
	if err != nil {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			return nil, nil, closeErr
		}
		return nil, nil, err
	}

	return inventoryConn, paymentConn, nil
}

func exitWithError(cancel context.CancelFunc, msg string, err error) {
	slog.Error(msg, "error", err)
	cancel()
	os.Exit(1)
}
