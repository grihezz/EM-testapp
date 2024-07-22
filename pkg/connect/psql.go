package connect

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func NewPostgresConnection(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		fmt.Println("Error parsing database config:", err)
		return nil, err
	}

	// Create a new connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Println("Error creating connection pool:", err)
		return nil, err
	}

	// Wait for a few seconds to ensure the connection is established
	time.Sleep(3 * time.Second)

	// Ping the database to check the connection
	err = pool.Ping(context.Background())
	if err != nil {
		fmt.Println("Error pinging the database:", err)
		return nil, err
	}

	fmt.Println("DB connection opened")

	return pool, nil
}
