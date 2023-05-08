package app

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/shamank/wb-l0/app/config"
	"github.com/shamank/wb-l0/app/internal/adapters/http"
	"github.com/shamank/wb-l0/app/internal/adapters/nats"
	"github.com/shamank/wb-l0/app/internal/repository"
	"github.com/shamank/wb-l0/app/internal/service"
	"github.com/shamank/wb-l0/app/server"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configDir string) {

	cfg, err := config.ConfigInit(configDir)
	if err != nil {
		logrus.Fatalf("error occurred in initial config: %s", err.Error())
		return
	}

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SSLMode))
	if err != nil {
		logrus.Fatalf("error occurred in connecting to postgres: %s", err.Error())
	}

	repos := repository.NewRepository(db)

	//cache := cache.NewFrom(cfg.Cache.TTL, cfg.Cache.CleanUp, repos.Orders.GetAll())

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	memcache := cache.New(cfg.Cache.TTL, cfg.Cache.CleanUp)

	CacheInit(ctx, repos, memcache)

	services := service.NewService(repos, memcache)

	handlers := http.NewHandler(services)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.InitAPI())

	sc, err := nats.NewNats(cfg.Nats, services, ctx)

	sc.InitSubscriptions()

	go func() {
		if err := srv.Start(); err != nil {
			return
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		return
	}
}

func CacheInit(ctx context.Context, repos *repository.Repository, c *cache.Cache) {
	items, err := repos.Orders.GetAll(ctx)
	if err != nil {
		logrus.Fatalf("error occurred in cache init")
		return
	}
	for _, item := range items {
		c.SetDefault(item.OrderUID, item)
	}
	logrus.Info("cache was loaded from postgres")

}
