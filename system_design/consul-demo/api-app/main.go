package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// consulCatalogService — ответ Consul API /v1/catalog/service/{name}
type consulCatalogService struct {
	ServiceAddress string `json:"ServiceAddress"`
	ServicePort    int    `json:"ServicePort"`
}

func main() {
	serviceName := getEnv("SERVICE_NAME", "api")
	port := getEnv("SERVICE_PORT", "8081")
	consulAddr := getEnv("CONSUL_HTTP_ADDR", "localhost:8500")

	// /call — проксирует запрос к web-app, адрес берет из Consul
	http.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		// 1. Service Query: спрашиваем Consul, где живет "web"
		addr, err := discoverService("web", consulAddr)
		if err != nil {
			http.Error(w, "service discovery failed: "+err.Error(), 502)
			return
		}

		// 2. Делаем запрос к найденному инстансу
		resp, err := http.Get("http://" + addr)
		if err != nil {
			http.Error(w, "upstream error: "+err.Error(), 502)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"from":"%s","upstream_addr":"%s","upstream_response":%s}`,
			serviceName, addr, string(body))
	})

	// /health — для Consul health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"passing"}`)
	})

	// /services — список всех сервисов в Consul
	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		services, err := listServices(consulAddr)
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(services)
	})

	log.Printf("[%s] listening on :%s", serviceName, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// discoverService — Service Query через Consul HTTP API
// Возвращает "host:port" для запроса к upstream
func discoverService(name, consulAddr string) (string, error) {
	url := fmt.Sprintf("http://%s/v1/health/service/%s?passing=true", consulAddr, name)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Используем health endpoint — возвращает только healthy инстансы
	var result []struct {
		Service struct {
			Address string `json:"Address"`
			Port    int    `json:"Port"`
		} `json:"Service"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result) == 0 {
		return "", fmt.Errorf("no healthy instances of '%s' found", name)
	}

	// Берем первый healthy инстанс
	svc := result[0].Service
	return fmt.Sprintf("%s:%d", svc.Address, svc.Port), nil
}

// listServices — список всех зарегистрированных сервисов
func listServices(consulAddr string) (map[string]interface{}, error) {
	url := fmt.Sprintf("http://%s/v1/catalog/services", consulAddr)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var _ = time.Now // suppress unused import
