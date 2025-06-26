package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/Ricky004/watchdata/pkg/otelpipeline/builderconfig"
// 	"github.com/Ricky004/watchdata/pkg/otelpipeline/exporter"
// 	"github.com/Ricky004/watchdata/pkg/otelpipeline/processors"
// 	"github.com/Ricky004/watchdata/pkg/otelpipeline/receviers"
// 	"github.com/Ricky004/watchdata/pkg/types/otelpipelinetypes"
// 	"gopkg.in/yaml.v3"
// )

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/Ricky004/watchdata/pkg/clickhousestore"
// )

// func main() {
// 	ctx := context.Background()

// 	config := clickhousestore.Config{
// 		Connection: clickhousestore.ConnectionConfig{
// 			DialTimeout: 5 * time.Second,
// 		},
// 	}

// 	_, err := clickhousestore.NewClickHouseProvider(ctx, config)
// 	if err != nil {
// 		log.Fatalf("error: %v", err)
// 	}

// 	fmt.Println("✅ ClickHouse connected")

// 	runTestQuery()

// }

// func runTestQuery() {
// 	dsn := "clickhouse://default:pass@localhost:9000/default"
// 	db, err := sql.Open("clickhouse", dsn)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to open DB: %v", err)
// 	}
// 	defer db.Close()

// 	if err := db.Ping(); err != nil {
// 		log.Fatalf("❌ Ping failed: %v", err)
// 	}

// 	// 🧱 Create a test table
// 	_, err = db.Exec(`
// 		CREATE TABLE IF NOT EXISTS users (
// 			id UInt32,
// 			name String,
// 			created_at DateTime
// 		) ENGINE = MergeTree()
// 		ORDER BY id
// 	`)
// 	if err != nil {
// 		log.Fatalf("❌ Create table failed: %v", err)
// 	}
// 	fmt.Println("🪵 Table 'users' created.")

// 	// ➕ Insert data
// 	_, err = db.Exec(`INSERT INTO users (id, name, created_at) VALUES (?, ?, ?)`, 1, "Alice", time.Now())
// 	if err != nil {
// 		log.Fatalf("❌ Insert failed: %v", err)
// 	}
// 	fmt.Println("✅ Data inserted.")

// 	// 🔍 Read data
// 	var (
// 		id        uint32
// 		name      string
// 		createdAt time.Time
// 	)
// 	err = db.QueryRow(`SELECT id, name, created_at FROM users LIMIT 1`).Scan(&id, &name, &createdAt)
// 	if err != nil {
// 		log.Fatalf("❌ Select failed: %v", err)
// 	}

// 	fmt.Printf("📦 Retrieved: ID=%d, Name=%s, CreatedAt=%s\n", id, name, createdAt.Format(time.RFC3339))
// }

// func main() {
// 	// Step 1: User selection (you'd get this from UI)
// 	selectedReceivers1 := []string{"otlp", "filelog"}
// 	selectedProcessors1 := []string{"batch"}
// 	selectedExporter1 := []string{"watchdataexporter"}

// 	selectedReceivers2 := []string{"filelog"}

// 	// Step 2: Build component blocks
// 	receivers := receviers.BuildReceivers(selectedReceivers1)
// 	processors := processors.BuildProcessors(selectedProcessors1)
// 	exporter := exporter.BuildExporter(selectedExporter1)

// 	// Step 3: Build pipeline section
// 	pipelines := map[string]otelpipelinetypes.Pipeline{
// 		"logs": {
// 			Receivers:  selectedReceivers1,
// 			Processors: selectedProcessors1,
// 			Exporters: selectedExporter1,
// 		},
// 		"logs/2": {
// 			Receivers:  selectedReceivers2,
// 		},
// 	}

// 	// Step 4: Wrap into final config
// 	config := otelpipelinetypes.OTelConfig{
// 		Receivers:  receivers,
// 		Processors: processors,
// 		Exporters: exporter,
// 		Service: otelpipelinetypes.ServiceConfig{
// 			Pipelines: pipelines,
// 		},
// 	}

// 	// Step 5: Marshal and write to otel-config.yaml
// 	data, err := yaml.Marshal(config)
// 	if err != nil {
// 		fmt.Println("Failed to marshal YAML:", err)
// 		os.Exit(1)
// 	}

// 	if err := os.WriteFile("otel-collector-config.yaml", data, 0644); err != nil {
// 		fmt.Println("Failed to write file:", err)
// 		os.Exit(1)
// 	}

// 	err = builderconfig.SyncBuilderConfig("builder-config.yaml", selectedReceivers1, selectedProcessors1, selectedExporter1)
// 	if err != nil {
// 		log.Fatalf("failed to update builder-config.yaml: %v", err)
// 	}

// 	fmt.Println("otel-config.yaml generated successfully 🎉")
// }