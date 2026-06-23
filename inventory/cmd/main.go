package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/fvaiiii/microservices-course-boilerplate/inventory/internal/interceptor"
	"github.com/fvaiiii/microservices-course-boilerplate/inventory/pkg/app"
)

const (
	grpcAddress = ":50051"

	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbURI := os.Getenv("DB_URI")
	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		exitWithError(cancel, "создание пула соединений", err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		exitWithError(cancel, "проверка соединений с БД", err)
	}
	slog.Info("подключение к PostgreSQL установлено")

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		exitWithError(cancel, "создание transaction manager", err)
	}

	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", grpcAddress)
	if err != nil {
		exitWithError(cancel, "не удалось создать listener", err)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true,
		}),
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	)

	app.RegisterServices(grpcServer, pool, txManager)

	reflection.Register(grpcServer)

	go func() {
		slog.Info("запуск InventoryService", "адрес", grpcAddress)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("ошибка запуска сервера", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("остановка gRPC сервера")
	grpcServer.GracefulStop()
	slog.Info("сервер остановлен")
}

func exitWithError(cancel context.CancelFunc, msg string, err error) {
	slog.Error(msg, "error", err)
	cancel()
	os.Exit(1)
}
