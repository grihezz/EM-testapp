package connect

import (
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"os"
)

// NewClickhouseConnection - функция создающая подключение к Clickhouse
func NewClickhouseConnection() (*sql.DB, error) {
	clickhouseDriver := os.Getenv("CLICKHOUSE_DRIVER")
	clickhouseSrc := "clickhouse://clickhouse:9000?debug=true"

	connect, err := sql.Open(clickhouseDriver, clickhouseSrc)
	if err != nil {
		fmt.Println("Error connecting to ClickHouse:", err)
		return nil, err
	}

	return connect, nil
}
