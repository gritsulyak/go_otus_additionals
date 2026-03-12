package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func newKafkaProducer() (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"acks":               "all", // ждём подтверждения от всех ISR реплик
		"retries":            10,
		"enable.idempotence": true, // идемпотентный продюсер: дедупликация на стороне брокера
		"compression.type":   "snappy",
		"linger.ms":          5, // небольшая пауза для батчинга
		"batch.size":         65536,
	})
}

func runOutboxWorker(db *sql.DB, p *kafka.Producer, topic string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Горутина для обработки delivery reports
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("delivery failed: %v", ev.TopicPartition.Error)
				}
			}
		}
	}()

	for range ticker.C {
		if err := processOutboxBatch(db, p, topic); err != nil {
			log.Println("outbox batch error:", err)
		}
	}
}

type outboxRecord struct {
	ID          int64
	AggregateID string
	EventType   string
	Payload     []byte
}

func processOutboxBatch(db *sql.DB, p *kafka.Producer, topic string) error {
	ctx := context.Background()

	rows, err := db.QueryContext(ctx,
		`SELECT id, aggregate_id, event_type, payload
         FROM outbox
         WHERE processed_at IS NULL
         ORDER BY id
         LIMIT 100`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var batch []outboxRecord
	for rows.Next() {
		var r outboxRecord
		if err := rows.Scan(&r.ID, &r.AggregateID, &r.EventType, &r.Payload); err != nil {
			return err
		}
		batch = append(batch, r)
	}

	for _, rec := range batch {
		if err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Key:   []byte(rec.AggregateID),
			Value: rec.Payload,
			Headers: []kafka.Header{
				{Key: "event_type", Value: []byte(rec.EventType)},
			},
		}, nil); err != nil {
			log.Printf("produce error for id=%d: %v", rec.ID, err)
			continue
		}

		// Помечаем как обработанное (только после успешного produce)
		if _, err := db.ExecContext(ctx,
			`UPDATE outbox SET processed_at = now() WHERE id = $1`,
			rec.ID,
		); err != nil {
			log.Printf("mark processed error for id=%d: %v", rec.ID, err)
		}
	}

	// Ждём доставки всех сообщений в батче
	p.Flush(3000)
	return nil
}
