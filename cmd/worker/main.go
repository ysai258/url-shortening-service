package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/url-shortening-service/pkg/db"
)

// worker periodically deletes expired links from the database
func worker(ticker *time.Ticker, db *sql.DB, done chan bool) {
	for {
		select {
		case <-ticker.C:
			_, err := db.Exec("DELETE FROM links WHERE created < NOW() - INTERVAL 1 DAY")
			if err != nil {
				log.Fatalf("Failed to delete expired links: %v", err)
			}
		case <-done:
			return
		}
	}
}

func main() {
	db.InitDB()
	ticker := time.NewTicker(5 * time.Minute)
	done := make(chan bool)

	// start the worker
	go worker(ticker, db.Db, done)

	// create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// listen for interrupt signals and stop the worker gracefully
	go func() {
		<-quit
		log.Println("Shutting down worker...")
		ticker.Stop()
		done <- true
	}()

	// wait for the worker to stop
	<-done
	log.Println("Worker has shut down.")
}
