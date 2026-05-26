package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// consulServiceRegistration — payload для регистрации в Consul
type consulServiceRegistration struct {
	ID      string      `json:"ID"`
	Name    string      `json:"Name"`
	Port    int         `json:"Port"`
	Address string      `json:"Address"`
	Check   consulCheck `json:"Check"`
}

type consulCheck struct {
	HTTP     string `json:"HTTP"`
	Interval string `json:"Interval"`
	Timeout  string `json:"Timeout"`
}

func main() {
	serviceName := getEnv("SERVICE_NAME", "web")
	port := getEnv("SERVICE_PORT", "8080")
	consulAddr := getEnv("CONSUL_HTTP_ADDR", "localhost:8500")

	// Регистрируем сервис в Consul при старте
	go registerInConsul(serviceName, port, consulAddr)

	// HTTP хендлеры
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"service":"%s","status":"ok","time":"%s"}`,
			serviceName, time.Now().Format(time.RFC3339))
	})

	// Health check endpoint — Consul будет опрашивать его каждые 10 секунд
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"passing"}`)
	})

	log.Printf("[%s] listening on :%s", serviceName, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func registerInConsul(name, port, consulAddr string) {
	// Ждем, пока Consul поднимется
	time.Sleep(3 * time.Second)

	registration := consulServiceRegistration{
		ID:      name + "-1",
		Name:    name,
		Port:    8080,
		Address: name + "-app", // hostname контейнера в docker network
		Check: consulCheck{
			// Consul дергает этот URL для health check
			HTTP:     fmt.Sprintf("http://%s-app:%s/health", name, port),
			Interval: "10s",
			Timeout:  "2s",
		},
	}

	body, _ := json.Marshal(registration)
	url := fmt.Sprintf("http://%s/v1/agent/service/register", consulAddr)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		log.Printf("[%s] failed to build consul request: %v", name, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[%s] consul registration failed: %v", name, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("[%s] registered in consul, status: %d", name, resp.StatusCode)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
