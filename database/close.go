package database

import (
	"context"
	"log"
	"time"
)

func (db *DB) Close() {

	// 🔹 Close SQL DB
	if db.SQLDB != nil {
		if err := db.SQLDB.Close(); err != nil {
			log.Println("⚠️ Error closing SQL DB:", err)
		} else {
			log.Println("🛑 SQL DB closed")
		}
	}

	// 🔹 Close PGX Pool
	if db.PGX != nil {

		// Optional: wait for ongoing queries (graceful shutdown)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		done := make(chan struct{})

		go func() {
			db.PGX.Close()
			close(done)
		}()

		select {
		case <-done:
			log.Println("🛑 PGX pool closed")
		case <-ctx.Done():
			log.Println("⚠️ PGX pool close timeout")
		}
	}
}
