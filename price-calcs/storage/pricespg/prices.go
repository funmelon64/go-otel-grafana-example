package pricespg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
	semconv "go.opentelemetry.io/otel/semconv/v1.22.0"
	"math/rand"
	"time"
)

const DRIVERS_COUNT = 100

type Storage struct {
	db *sql.DB
}

func NewStorage(addr string, dbname string, user, pass string) (*Storage, error) {
	const connStrTemplate = "postgres://%s:%s@%s/%s?sslmode=disable"
	connStr := fmt.Sprintf(connStrTemplate, user, pass, addr, dbname)

	db, err := otelsql.Open("postgres", connStr, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	storage := &Storage{db: db}

	if err := storage.mustInitTable(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) mustInitTable() error {
	const pricesTabQuery = "CREATE TABLE IF NOT EXISTS prices (id SERIAL PRIMARY KEY, price DECIMAL(10, 2) NOT NULL, time TIMESTAMP NOT NULL, driver_id INT NOT NULL);"
	const discountsTabQuery = "CREATE TABLE IF NOT EXISTS discounts (id SERIAL PRIMARY KEY, discount INT NOT NULL, driver_id INT NOT NULL);"
	_, err := s.db.Exec(pricesTabQuery)
	if err != nil {
		return fmt.Errorf("failed to create prices table: %w", err)
	}
	_, err = s.db.Exec(discountsTabQuery)
	if err != nil {
		return fmt.Errorf("failed to create discounts table: %w", err)
	}

	// Вставка тестовых данных
	for i := 1; i <= DRIVERS_COUNT; i++ {
		price := rand.Float64() * 10000
		time := time.Now()

		_, err = s.db.Exec("INSERT INTO prices (price, time, driver_id) VALUES ($1, $2, $3)", price, time, i)
		if err != nil {
			return fmt.Errorf("failed to insert into prices table: %w", err)
		}

		for j := 0; j < 3; j++ {
			discount := rand.Intn(100)
			_, err = s.db.Exec("INSERT INTO discounts (discount, driver_id) VALUES ($1, $2)", discount, i)
			if err != nil {
				return fmt.Errorf("failed to insert into discounts table: %w", err)
			}
		}
	}

	return nil
}

func (s *Storage) GetDriverPrice(ctx context.Context, driverId string) (float64, error) {
	var id float64
	query := `SELECT price FROM prices WHERE driver_id = $1`

	// Важно передавать ctx в запрос, чтобы запрос был частью трейса
	err := s.db.QueryRowContext(ctx, query, driverId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetDriverDiscounts(ctx context.Context, driverId string) ([]int, error) {
	query := `SELECT discount FROM discounts WHERE driver_id = $1`

	// Важно передавать ctx в запрос, чтобы запрос был частью трейса
	rows, err := s.db.QueryContext(ctx, query, driverId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var discounts []int
	for rows.Next() {
		var discount int
		if err := rows.Scan(&discount); err != nil {
			return nil, err
		}
		discounts = append(discounts, discount)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return discounts, nil
}
