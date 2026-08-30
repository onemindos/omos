// om-streaming-bridge: NATS → RisingWave/Flink/Postgres
// Reads NATS JetStream, writes to streaming databases
// Deploy as NEX native workload

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	_ "github.com/lib/pq"
)

const (
	batchSize     = 100
	flushInterval = 5 * time.Second
)

type BridgeConfig struct {
	NATSURL       string
	RisingWaveURL string
	FlinkURL      string
	PostgresURL   string
}

type EventBatch struct {
	events []map[string]interface{}
	mu     chan struct{}
}

var (
	nc       *nats.Conn
	js       nats.JetStreamContext
	rwDB     *sql.DB
	flinkDB  *sql.DB
	pgDB     *sql.DB
	evtBatch *EventBatch
)

func main() {
	cfg := BridgeConfig{
		NATSURL:       getEnv("NATS_URL", "nats://127.0.0.1:4222"),
		RisingWaveURL: getEnv("RISINGWAVE_URL", ""),
		FlinkURL:      getEnv("FLINK_URL", ""),
		PostgresURL:   getEnv("POSTGRES_URL", ""),
	}

	// Connect to NATS
	var err error
	nc, err = nats.Connect(cfg.NATSURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NATS connect failed: %v\n", err)
		os.Exit(1)
	}
	defer nc.Close()

	js, err = nc.JetStream()
	if err != nil {
		fmt.Fprintf(os.Stderr, "JetStream error: %v\n", err)
		os.Exit(1)
	}

	// Connect to sinks
	if cfg.RisingWaveURL != "" {
		rwDB, err = sql.Open("postgres", cfg.RisingWaveURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "RisingWave connect failed: %v\n", err)
		}
		defer rwDB.Close()
		initRisingWaveTables()
	}

	if cfg.FlinkURL != "" {
		flinkDB, err = sql.Open("postgres", cfg.FlinkURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Flink connect failed: %v\n", err)
		}
		defer flinkDB.Close()
	}

	if cfg.PostgresURL != "" {
		pgDB, err = sql.Open("postgres", cfg.PostgresURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Postgres connect failed: %v\n", err)
		}
		defer pgDB.Close()
	}

	// Initialize batch buffer
	evtBatch = &EventBatch{
		events: make([]map[string]interface{}, 0, batchSize),
		mu:     make(chan struct{}, 1),
	}

	// Subscribe to streams
	subscribeStreams()

	// Start batch flusher
	go batchFlusher()

	fmt.Println("om-streaming-bridge running")
	fmt.Printf("Sinks: RisingWave=%v Flink=%v Postgres=%v\n",
		rwDB != nil, flinkDB != nil, pgDB != nil)

	select {}
}

func subscribeStreams() {
	streams := []struct {
		stream  string
		subject string
		table   string
	}{
		{"EVT", "EVT.>", "raw_events"},
		{"COT", "COT.>", "raw_cot"},
		{"VOICE", "VOICE.>", "raw_voice"},
		{"CMD", "CMD.>", "raw_commands"},
		{"TEL", "TEL.>", "raw_telemetry"},
	}

	for _, s := range streams {
		go func(stream, subject, table string) {
			sub, err := js.Subscribe(subject, func(msg *nats.Msg) {
				handleMessage(msg, table)
			}, nats.Durable("streaming-bridge-"+stream), nats.DeliverAll())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Subscribe %s error: %v\n", stream, err)
				return
			}
			_ = sub
			fmt.Printf("Subscribed to %s (%s) → %s\n", stream, subject, table)
		}(s.stream, s.subject, s.table)
	}
}

func handleMessage(msg *nats.Msg, table string) {
	var data map[string]interface{}
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		// Store raw if not JSON
		data = map[string]interface{}{
			"subject": msg.Subject,
			"data":    string(msg.Data),
			"time":    time.Now(),
		}
	}
	data["_subject"] = msg.Subject
	data["_time"] = time.Now()

	evtBatch.mu <- struct{}{}
	evtBatch.events = append(evtBatch.events, data)
	if len(evtBatch.events) >= batchSize {
		flushBatch(table)
	}
	<-evtBatch.mu
}

func batchFlusher() {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for range ticker.C {
		tables := []string{"raw_events", "raw_cot", "raw_voice", "raw_commands", "raw_telemetry"}
		for _, table := range tables {
			evtBatch.mu <- struct{}{}
			if len(evtBatch.events) > 0 {
				flushBatch(table)
			}
			<-evtBatch.mu
		}
	}
}

func flushBatch(table string) {
	if len(evtBatch.events) == 0 {
		return
	}

	batch := evtBatch.events
	evtBatch.events = make([]map[string]interface{}, 0, batchSize)

	// Insert to RisingWave
	if rwDB != nil {
		insertToRisingWave(batch, table)
	}

	// Insert to Flink
	if flinkDB != nil {
		insertToFlink(batch, table)
	}

	// Insert to Postgres
	if pgDB != nil {
		insertToPostgres(batch, table)
	}

	fmt.Printf("Flushed %d events to %s\n", len(batch), table)
}

func insertToRisingWave(batch []map[string]interface{}, table string) {
	// Bulk insert with prepared statement
	stmt, err := rwDB.Prepare(fmt.Sprintf("INSERT INTO %s (data) VALUES ($1)", table))
	if err != nil {
		// Table might not exist, create it
		createTable(rwDB, table)
		stmt, _ = rwDB.Prepare(fmt.Sprintf("INSERT INTO %s (data) VALUES ($1)", table))
	}
	defer stmt.Close()

	for _, event := range batch {
		data, _ := json.Marshal(event)
		_, _ = stmt.Exec(string(data))
	}
}

func insertToFlink(batch []map[string]interface{}, table string) {
	if flinkDB == nil {
		return
	}
	for _, event := range batch {
		data, _ := json.Marshal(event)
		_, _ = flinkDB.Exec(fmt.Sprintf("INSERT INTO %s VALUES (?)", table), string(data))
	}
}

func insertToPostgres(batch []map[string]interface{}, table string) {
	if pgDB == nil {
		return
	}
	for _, event := range batch {
		data, _ := json.Marshal(event)
		_, _ = pgDB.Exec(fmt.Sprintf("INSERT INTO %s (data, time) VALUES ($1, NOW())", table), string(data))
	}
}

func initRisingWaveTables() {
	tables := []string{"raw_events", "raw_cot", "raw_voice", "raw_commands", "raw_telemetry"}
	for _, table := range tables {
		createTable(rwDB, table)
	}
}

func createTable(db *sql.DB, table string) {
	if db == nil {
		return
	}
	// RisingWave/Postgres JSON table
	_, _ = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		data JSONB,
		time TIMESTAMPTZ DEFAULT NOW(),
		subject TEXT
	)`, table))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
