package bookingpg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
	semconv "go.opentelemetry.io/otel/semconv/v1.22.0"
	"time"
)

// Booking представляет собой запись о бронировании
type Booking struct {
	ID    int
	Price float64
	Time  time.Time
}

// Storage предоставляет методы для работы с базой данных
type Storage struct {
	db *sql.DB
}

// NewStorage создает новое соединение с базой данных
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
	query := "CREATE TABLE IF NOT EXISTS bookings (id SERIAL PRIMARY KEY, price DECIMAL(10, 2) NOT NULL, time TIMESTAMP NOT NULL);"
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

// AddBooking добавляет новое бронирование в базу данных
func (s *Storage) AddBooking(ctx context.Context, price float64, time time.Time) (int, error) {
	var id int
	query := `INSERT INTO bookings (price, time) VALUES ($1, $2) RETURNING id`
	err := s.db.QueryRowContext(ctx, query, price, time).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetBookingById получает бронирование по ID
func (s *Storage) GetBookingById(ctx context.Context, id int) (*Booking, error) {
	var booking Booking
	query := `SELECT id, price, time FROM bookings WHERE id = $1`
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&booking.ID, &booking.Price, &booking.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no booking with id %d", id)
		}
		return nil, err
	}
	return &booking, nil
}
