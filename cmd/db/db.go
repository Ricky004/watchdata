package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Ricky004/watchdata/pkg/clickhousestore"
)

func main() {
	ctx := context.Background()

	config := clickhousestore.Config{
		Connection: clickhousestore.ConnectionConfig{
			DialTimeout: 5 * time.Second,
		},
	}

	_, err := clickhousestore.NewClickHouseProvider(ctx, config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("✅ ClickHouse connected")

	runTestQuery()

}

func runTestQuery() {
	dsn := "clickhouse://default:pass@localhost:9000/default"
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Ping failed: %v", err)
	}

	// 🧱 Create a test table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UInt32,
			name String,
			created_at DateTime
		) ENGINE = MergeTree()
		ORDER BY id
	`)
	if err != nil {
		log.Fatalf("❌ Create table failed: %v", err)
	}
	fmt.Println("🪵 Table 'users' created.")

	// ➕ Insert data
	_, err = db.Exec(`INSERT INTO users (id, name, created_at) VALUES (?, ?, ?)`, 1, "Alice", time.Now())
	if err != nil {
		log.Fatalf("❌ Insert failed: %v", err)
	}
	fmt.Println("✅ Data inserted.")

	// 🔍 Read data
	var (
		id        uint32
		name      string
		createdAt time.Time
	)
	err = db.QueryRow(`SELECT id, name, created_at FROM users LIMIT 1`).Scan(&id, &name, &createdAt)
	if err != nil {
		log.Fatalf("❌ Select failed: %v", err)
	}

	fmt.Printf("📦 Retrieved: ID=%d, Name=%s, CreatedAt=%s\n", id, name, createdAt.Format(time.RFC3339))
}
