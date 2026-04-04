package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rajesh_bond/production/config"
)

type DB struct {
	SQLDB *sql.DB
	PGX   *pgxpool.Pool
}

// NewDB initializes both sql.DB and pgxpool
func NewDB(cfg *config.Config) *DB {

	// 🔹 DSN for sql.DB
	sqlDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHOST,
		cfg.DBPORT,
		cfg.DBUSER,
		cfg.DBPASS,
		cfg.DBNAME,
		cfg.DBSSL,
	)

	// 🔹 DSN for pgx
	pgxDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUSER,
		cfg.DBPASS,
		cfg.DBHOST,
		cfg.DBPORT,
		cfg.DBNAME,
		cfg.DBSSL,
	)

	// ✅ Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// =========================
	// 🔹 1. Setup sql.DB
	// =========================
	sqlDB, err := sql.Open("postgres", sqlDSN)
	if err != nil {
		log.Fatal("❌ Failed to open SQL DB:", err)
	}

	// 🔥 Connection pool tuning (IMPORTANT)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// ✅ Ping check
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatal("❌ SQL DB ping failed:", err)
	}

	log.Println("✅ SQL DB connected successfully")

	// =========================
	// 🔹 2. Setup pgxpool
	// =========================
	pgxConfig, err := pgxpool.ParseConfig(pgxDSN)
	if err != nil {
		log.Fatal("❌ Failed to parse PGX config:", err)
	}

	// 🔥 Pool tuning (IMPORTANT)
	pgxConfig.MaxConns = 10
	pgxConfig.MinConns = 2
	pgxConfig.MaxConnLifetime = time.Hour
	pgxConfig.MaxConnIdleTime = 30 * time.Minute

	pgxPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatal("❌ Failed to connect PGX:", err)
	}

	// ✅ Ping check
	if err := pgxPool.Ping(ctx); err != nil {
		log.Fatal("❌ PGX ping failed:", err)
	}

	log.Println("✅ PGX connected successfully")

	// =========================
	// 🔹 Return combined DB
	// =========================
	return &DB{
		SQLDB: sqlDB,
		PGX:   pgxPool,
	}
}
