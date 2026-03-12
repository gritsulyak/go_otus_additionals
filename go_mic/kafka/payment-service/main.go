package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Payment struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	if err := runMigrations(db); err != nil {
		log.Fatal("migration error:", err)
	}

	// Запускаем outbox-воркер в фоне
	go func() {
		p, err := newKafkaProducer()
		if err != nil {
			log.Fatal("producer error:", err)
		}
		defer p.Close()
		runOutboxWorker(db, p, "payments.created",
			time.Duration(500)*time.Millisecond)
	}()

	http.HandleFunc("/payments", handleCreatePayment)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("payment-service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var p Payment
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if p.ID == "" || p.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 1. Бизнес-запись
	_, err = tx.ExecContext(ctx,
		`INSERT INTO payments (id, amount) VALUES ($1, $2)
         ON CONFLICT (id) DO NOTHING`,
		p.ID, p.Amount,
	)
	if err != nil {
		log.Println("insert payment error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 2. Outbox-запись в той же транзакции
	payload, _ := json.Marshal(map[string]any{
		"id":     p.ID,
		"amount": p.Amount,
		"ts":     time.Now().UTC(),
	})
	_, err = tx.ExecContext(ctx,
		`INSERT INTO outbox (aggregate_id, event_type, payload)
         VALUES ($1, $2, $3)`,
		p.ID, "PaymentCreated", payload,
	)
	if err != nil {
		log.Println("insert outbox error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("commit error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func runMigrations(db *sql.DB) error {
	data, err := os.ReadFile("migrations/init.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(data))
	return err
}
