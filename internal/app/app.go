package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/begenov/backend/internal/config"
	delivery "github.com/begenov/backend/internal/delivery/http"
	"github.com/begenov/backend/internal/repository"
	"github.com/begenov/backend/internal/server"
	"github.com/begenov/backend/internal/service"
	"github.com/begenov/backend/pkg/auth"
	"github.com/begenov/backend/pkg/db"
	"github.com/begenov/backend/pkg/hash"
)

func Run(cfg *config.Config) error {
	db, err := db.NewDB(cfg.Postgres.Driver, cfg.Postgres.DSN)
	if err != nil {
		return err
	}

	hash := hash.NewHash()

	token, err := auth.NewManager(cfg.JWT.TokenSymmetricKey)
	if err != nil {
		return err
	}

	repo := repository.NewRepository(db)

	service := service.NewService(repo, hash, token, cfg.JWT.AccessTokenDuration)

	handler := delivery.NewHandler(service, token)

	srv := server.NewServer(cfg, handler.Init(cfg))

	go func() {
		if err = srv.Run(); err != nil {
			log.Fatalf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	log.Println("Server started")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Printf("failed to stop server: %v", err)
	}

	return nil
}
