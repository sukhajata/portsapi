package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// TimescaleEngine implements SQLEngine
type PostgresEngine struct {
	pool *pgxpool.Pool
}

// FatalError implements error, indicates connection lost, action needs to be taken
type FatalError struct {
	message string
}

func (e *FatalError) Error() string {
	return e.message
}

// NewPostgresEngine factory method for creating postgres engine
func NewPostgresEngine(psqlURL string) (*PostgresEngine, error) {
	var err error
	var pool *pgxpool.Pool
	retries := 0
	for {
		pool, err = pgxpool.Connect(context.Background(), psqlURL)
		if err != nil {
			retries++
			if retries > 5 {
				return nil, err
			}
			time.Sleep(time.Second * 2)
			continue
		}

		break
	}

	return &PostgresEngine{
		pool: pool,
	}, nil

}

// Query - get array
func (t *PostgresEngine) Query(queryString string, arguments ...interface{}) ([]interface{}, error) {
	conn, err := t.pool.Acquire(context.Background())
	if err != nil {
		return nil, &FatalError{message: err.Error()}
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), queryString, arguments...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]interface{}, 0)

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return results, err
		}
		results = append(results, values)
	}

	return results, err
}

// Exec - run a query without return
func (t *PostgresEngine) Exec(queryString string, arguments ...interface{}) error {
	conn, err := t.pool.Acquire(context.Background())
	if err != nil {
		return &FatalError{message: err.Error()}
	}

	defer conn.Release()
	_, err = conn.Exec(context.Background(), queryString, arguments...)

	return err
}

// ScanRow - query a row and scan into the value pointer
func (t *PostgresEngine) ScanRow(queryString string, valuePtr interface{}, arguments ...interface{}) error {
	conn, err := t.pool.Acquire(context.Background())
	if err != nil {
		return &FatalError{message: err.Error()}
	}

	defer conn.Release()
	err = conn.QueryRow(context.Background(), queryString, arguments...).Scan(valuePtr)

	return err
}

// Close the pool
func (t *PostgresEngine) Close() {
	// check pool is open
	conn, err := t.pool.Acquire(context.Background())
	if err == nil {
		// no error so it must be open
		conn.Release()
		t.pool.Close()
	}
}
