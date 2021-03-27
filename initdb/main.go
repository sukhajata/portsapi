package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	dbPool    *pgxpool.Pool
	adminPool *pgxpool.Pool
	adminURL  = os.Getenv("adminURL")
	dbURL     = os.Getenv("dbURL")
	dbName    = os.Getenv("dbName")
)

func main() {
	var err error
	tries := 0
	for {
		//connect
		adminPool, err = pgxpool.Connect(context.Background(), adminURL)
		if err != nil {
			tries++
			if tries > 5 {
				panic(err)
			}
			fmt.Println(err)
			time.Sleep(time.Second * 2)
			continue
		}
		break
	}

	fmt.Println("Connected to postgres")

	sql := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname=$1")
	fmt.Println(sql)
	rows, err := adminPool.Query(context.Background(), sql, dbName)
	if err != nil {
		panic(err)
	}

	if !rows.Next() {
		sql = fmt.Sprintf("CREATE DATABASE %s", dbName)
		fmt.Println(sql)
		_, err = adminPool.Exec(context.Background(), sql)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created database %s\n", dbName)
	} else {
		fmt.Printf("Not creating %s, database already exists\n", dbName)
	}

	adminPool.Close()

	dbPool, err = pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}

	sql = `CREATE TABLE IF NOT EXISTS ports (
    	id TEXT PRIMARY KEY,
    	data JSONB
  	)`
	_, err = dbPool.Exec(context.Background(), sql)
	if err != nil {
		panic(err)
	}

	dbPool.Close()
}
