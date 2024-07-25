package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"tt/internal/adapters/repository"
	"tt/internal/handlers/router"
	"tt/internal/service"
	"tt/pkg/connect"
	"tt/pkg/logger"
	"tt/pkg/migrations"
)

func main() {
	log := setupLogger(os.Getenv("ENV"))
	log.Info(
		"Starting application",
		slog.String("env", os.Getenv("ENV")),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	postgreConn, err := connect.NewPostgresConnection(os.Getenv("DSN"))
	if err != nil {
		log.Error("Connecting to SQL database error: ", err)
		return
	}
	defer postgreConn.Close()
	fmt.Println("DB connection opened")

	err = migrations.UpMigration(context.Background(), postgreConn)
	if err != nil {
		log.Error("Failed to up migration: ", err)
	}
	fmt.Println("Migration up success")

	clickhouse, err := connect.NewClickhouseConnection()
	if err != nil {
		log.Error("Main NewClickhouseConnection Error: ", err)
	}
	defer clickhouse.Close()
	fmt.Println("Clickhouse connection opened")

	err = migrations.UpClickhouse(context.Background(), clickhouse)
	fmt.Println("Migration up success сlickhouse")

	userRepository := repository.NewUserRepository(postgreConn)
	userService := service.NewUserService(userRepository)

	taskRepository := repository.NewTaskRepository(postgreConn)
	taskService := service.NewTaskService(taskRepository)

	r := router.SetupRouter(userService, taskService, log)

	log.Info("Starting client on port " + os.Getenv("PORT"))

	srv := &http.Server{
		Addr:    "0.0.0.0:" + os.Getenv("PORT"),
		Handler: r,
	}

	log.Info("starting server",
		"type", "START",
		"addr", os.Getenv("PORT"),
	)

	err = srv.ListenAndServe()
	if err != nil {
		log.Error("Server dropped: ", err)
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = setupPrettySlog()
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := logger.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
