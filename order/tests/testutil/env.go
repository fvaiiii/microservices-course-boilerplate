package testutil

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	invApp "github.com/fvaiiii/microservices-course-boilerplate/inventory/pkg/app"
	"github.com/fvaiiii/microservices-course-boilerplate/order/pkg/app"
	payApp "github.com/fvaiiii/microservices-course-boilerplate/payment/pkg/app"
	inventoryv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/fvaiiii/microservices-course-boilerplate/shared/pkg/proto/payment/v1"
)

const bufSize = 1024 * 1024

// Env — изолированное тестовое окружение: свои БД, свои сервисы.
// Каждый параллельный тест получает свой Env и не пересекается с другими.
type Env struct {
	HTTPClient *http.Client
	BaseURL    string

	InventoryClient inventoryv1.InventoryServiceClient
	PaymentClient   paymentv1.PaymentServiceClient

	// Пулы прямого доступа к БД для проверок состояния и seed-данных.
	OrderPool     *pgxpool.Pool
	InventoryPool *pgxpool.Pool

	// Имена изолированных БД (полезно для отладки).
	OrderDBName     string
	InventoryDBName string
}

// NewEnv поднимает окружение для одного теста и регистрирует cleanup.
func NewEnv(t *testing.T) *Env {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	orderDB := createIsolatedDB(ctx, t, "order", "../../migrations/order")
	t.Cleanup(orderDB.cleanup)

	inventoryDB := createIsolatedDB(ctx, t, "inventory", "../../migrations/inventory")
	t.Cleanup(inventoryDB.cleanup)

	orderPool := newTestPool(ctx, t, "order", orderDB.DSN)
	inventoryPool := newTestPool(ctx, t, "inventory", inventoryDB.DSN)

	txManager, err := manager.New(trmpgx.NewDefaultFactory(orderPool))
	if err != nil {
		t.Fatalf("txManager: %v", err)
	}

	invTxManager, err := manager.New(trmpgx.NewDefaultFactory(inventoryPool))
	if err != nil {
		t.Fatalf("invTxManager: %v", err)
	}

	// Inventory gRPC через bufconn.
	invLis := bufconn.Listen(bufSize)
	invServer := grpc.NewServer(invApp.Interceptors()...)
	invApp.RegisterServices(invServer, inventoryPool, invTxManager)
	go func() {
		if err := invServer.Serve(invLis); err != nil {
			t.Logf("invServer: %v", err)
		}
	}()
	t.Cleanup(invServer.Stop)

	invConn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return invLis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("invConn: %v", err)
	}
	t.Cleanup(func() {
		if err := invConn.Close(); err != nil {
			t.Logf("close invConn: %v", err)
		}
	})
	invClient := inventoryv1.NewInventoryServiceClient(invConn)

	// Payment gRPC через bufconn.
	payLis := bufconn.Listen(bufSize)
	payServer := grpc.NewServer(payApp.Interceptors()...)
	payApp.RegisterServices(payServer)
	go func() {
		if err := payServer.Serve(payLis); err != nil {
			t.Logf("payServer: %v", err)
		}
	}()
	t.Cleanup(payServer.Stop)

	payConn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return payLis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("payConn: %v", err)
	}
	t.Cleanup(func() {
		if err := payConn.Close(); err != nil {
			t.Logf("close payConn: %v", err)
		}
	})
	payClient := paymentv1.NewPaymentServiceClient(payConn)

	// Order HTTP через httptest.
	orderHandler, err := app.NewHTTPHandler(orderPool, txManager, invClient, payClient)
	if err != nil {
		t.Fatalf("order handler: %v", err)
	}
	ts := httptest.NewServer(orderHandler)
	t.Cleanup(ts.Close)

	return &Env{
		HTTPClient:      &http.Client{Timeout: 10 * time.Second},
		BaseURL:         ts.URL,
		InventoryClient: invClient,
		PaymentClient:   payClient,
		OrderPool:       orderPool,
		InventoryPool:   inventoryPool,
		OrderDBName:     orderDB.Name,
		InventoryDBName: inventoryDB.Name,
	}
}

func newTestPool(ctx context.Context, t *testing.T, name, dsn string) *pgxpool.Pool {
	t.Helper()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("%s pool: %v", name, err)
	}
	t.Cleanup(pool.Close)

	return pool
}
