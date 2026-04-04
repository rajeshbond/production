package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/rajesh_bond/production/config"
)

func NewDB(cfg *config.Config) *sql.DB {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHOST,
		cfg.DBPORT,
		cfg.DBUSER,
		cfg.DBPASS,
		cfg.DBNAME,
		cfg.DBSSL,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	// ✅ Create timeout context (5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ✅ Use PingContext instead of Ping
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	log.Println("✅ Database connected successfully")

	return db
}
