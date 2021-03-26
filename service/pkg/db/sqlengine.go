package db

// SQLEngine provides a wrapper over a sql db library
type SQLEngine interface {
	// Exec - execute a query with no result rows
	Exec(sql string, arguments ...interface{}) error

	// Query - select rows
	Query(sql string, arguments ...interface{}) ([]interface{}, error)

	// ScanRow - retrieve a single row into valuePtr
	ScanRow(sql string, valuePtr interface{}, arguments ...interface{}) error

	// Close - close the db connection
	Close()
}
