package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	_ "github.com/lib/pq"
)

type PaymentEvent struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"`
	Ts     string `json:"ts"`
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(20)
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatal("migration error:", err)
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":              os.Getenv("KAFKA_GROUP_ID"),
		"auto.offset.reset":     "earliest",
		"enable.auto.commit":    false, // ручной коммит после обработки
		"max.poll.interval.ms":  300000,
		"session.timeout.ms":    30000,
		"heartbeat.interval.ms": 3000,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if err := c.SubscribeTopics([]string{"payments.created"}, nil); err != nil {
		log.Fatal(err)
	}

	log.Println("billing-service consuming from payments.created")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigCh:
			log.Println("shutting down gracefully")
			return
		default:
		}

		msg, err := c.ReadMessage(200 * time.Millisecond)
		if err != nil {
			if e, ok := err.(kafka.Error); ok && e.Code() == kafka.ErrTimedOut {
				continue
			}
			log.Println("read error:", err)
			continue
		}

		if err := handleWithInbox(db, msg); err != nil {
			log.Printf("handle error for key=%s: %v", msg.Key, err)
			// Не коммитим — сообщение будет переполучено
			continue
		}

		// Ручной коммит только после успешной обработки
		if _, err := c.CommitMessage(msg); err != nil {
			log.Println("commit offset error:", err)
		}
	}
}

func handleWithInbox(db *sql.DB, msg *kafka.Message) error {
	messageID := string(msg.Key)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Inbox-проверка: атомарная INSERT с проверкой дубликата
	var inserted bool
	err = tx.QueryRowContext(ctx,
		`INSERT INTO inbox (message_id)
         VALUES ($1)
         ON CONFLICT (message_id) DO NOTHING
         RETURNING true`,
		messageID,
	).Scan(&inserted)

	if err == sql.ErrNoRows {
		// Дубликат — сообщение уже обрабатывалось
		log.Printf("duplicate message skipped: %s", messageID)
		return nil
	}
	if err != nil {
		return err
	}

	var evt PaymentEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO billing (id, amount) VALUES ($1, $2)
         ON CONFLICT (id) DO NOTHING`,
		evt.ID, evt.Amount,
	); err != nil {
		return err
	}

	log.Printf("billing created for payment_id=%s amount=%d", evt.ID, evt.Amount)
	return tx.Commit()
}

func runMigrations(db *sql.DB) error {
	data, err := os.ReadFile("migrations/init.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(data))
	return err
}
