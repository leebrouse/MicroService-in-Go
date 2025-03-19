package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/leebrouse/MicroService-in-Go/authentication-service/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var count int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

	log.Println("Starting authentication service")
	//TODO: connect the database
	conn := connectToDB()
	if conn == nil {
		log.Panicln("Can`t connect to the Postgres")
	}
	//set config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.router(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// open db
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close() // 避免泄漏
		return nil, err
	}

	return db, nil
}

// connect to the pg
func connectToDB() *sql.DB {
	//using os package to grab the dsn
	dsn := os.Getenv("DSN")
	fmt.Println(dsn)
	//loop and the max count can`t surpass 3 times
	for {
		//call openDB
		connection, err := openDB(dsn)
		if err != nil {
			fmt.Println("Postgres not ready yet")
			count++
		} else {
			log.Println("Connect to the Postgres")
			return connection
		}

		//surpass 10 time then break
		if count > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}

}
