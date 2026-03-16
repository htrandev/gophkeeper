package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/htrandev/gophkeeper/internal/authorizer"
	"github.com/htrandev/gophkeeper/internal/config"
	"github.com/htrandev/gophkeeper/internal/server"
	"github.com/htrandev/gophkeeper/internal/service"
	"github.com/htrandev/gophkeeper/internal/storage/postgres"
	"github.com/htrandev/gophkeeper/internal/storage/postgres/migrations"
	"github.com/htrandev/gophkeeper/pkg/logger"
	pb "github.com/htrandev/gophkeeper/proto"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("run ends with error: %s", err.Error())
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	log.Println("init config")
	cfg, err := config.InitServerConfig()
	if err != nil {
		return fmt.Errorf("init config: %w", err)
	}

	log.Println("init logger")
	zl, err := logger.NewZapLogger(cfg.LogLvl)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	zl.Info("init db")
	// zl.Info("", zap.String("dsn", cfg.DatabaseDsn))
	db, err := sql.Open("pgx", cfg.DatabaseDsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	zl.Info("ping db")
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	zl.Info("init storage")
	storage := postgres.New(cfg.MaxRetry, db)

	zl.Info("init provider")
	provider, err := goose.NewProvider(database.DialectPostgres, db, migrations.Embed)
	if err != nil {
		return fmt.Errorf("goose: create new provider: %w", err)
	}

	zl.Info("up migrations")
	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("goose: provider up: %w", err)
	}

	zl.Info("init grpc listener")
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return fmt.Errorf("init grpc listener: %w", err)
	}

	zl.Info("init grpc server")
	grpcSrv := grpc.NewServer()
	defer grpcSrv.GracefulStop()

	zl.Info("init authorizer")
	auth := authorizer.New(cfg.Signature, cfg.TokenTTL)

	zl.Info("init auth service")
	authService := service.NewAuth(auth, storage)

	zl.Info("init keeper service")
	keeperService := service.NewKeeper(auth, storage)

	zl.Info("register auth server")
	pb.RegisterAuthServer(grpcSrv, server.NewAuth(
		server.WithLogger(zl),
		server.WithAuthService(authService),
	))
	zl.Info("register gophkeeper server")
	pb.RegisterGophkeeperServer(grpcSrv, server.NewGophkeeper(
		server.WithLogger(zl),
		server.WithAuthorizerService(auth),
		server.WithKeeperService(keeperService),
	))

	errc := make(chan error, 1)
	go func() {
		zl.Info("start serving grpc", zap.String("addr", cfg.Addr))
		if err := grpcSrv.Serve(lis); err != nil && err != http.ErrServerClosed {
			err = fmt.Errorf("serve server: %v", err)
			errc <- err
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-errc:
		if err != nil {
			return fmt.Errorf("serve ends with err: %w", err)
		}
	}

	return nil
}
