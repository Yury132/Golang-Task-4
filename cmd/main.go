package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Yury132/Golang-Task-4/internal/client/api"
	"github.com/Yury132/Golang-Task-4/internal/config"
	"github.com/Yury132/Golang-Task-4/internal/service"
	"github.com/Yury132/Golang-Task-4/internal/storage"
	transport "github.com/Yury132/Golang-Task-4/internal/transport/http"
	"github.com/Yury132/Golang-Task-4/internal/transport/http/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	dialect        = "pgx"
	commandUp      = "up"
	commandDown    = "down"
	migrationsPath = "./internal/migrations"
)

func main() {

	// Конфигурации
	cfg, err := config.Parse()
	if err != nil {
		panic(err)
	}

	// Логгер
	logger := cfg.Logger()

	// Миграции
	db, err := goose.OpenDBWithDriver(dialect, cfg.GetDBConnString())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open db by goose")
	}

	if err = goose.Run(commandUp, db, migrationsPath); err != nil {
		logger.Fatal().Msgf("migrate %v: %v", commandUp, err)
	}

	if err = db.Close(); err != nil {
		logger.Fatal().Err(err).Msg("failed to close db connection by goose")
	}

	// Настройка БД
	poolCfg, err := cfg.PgPoolConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to DB")
	}

	// Подключение к БД
	conn, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to db")
	}

	// Для обогащения сообщений
	userAPI := api.New(logger)

	// Хранилище
	strg := storage.New(conn)
	// Сервис
	svc := service.New(logger, userAPI, strg)
	// Хэндлер
	handler := handlers.New(logger, svc)
	// Сервер
	srv := transport.New(":8080").WithHandler(handler)

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)

	// Запусвкаем сервер
	go func() {
		fmt.Println("Сервер запущен")
		if err = srv.Run(); err != nil {
			logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// Ждем нажатия Ctrl+C
	<-shutdown
}
