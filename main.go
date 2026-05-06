package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type configJSON struct {
	N8N_API_URL string `json:"N8N_API_URL"`
	N8N_API_KEY string `json:"N8N_API_KEY"`
}

func main() {
	configStr := flag.String("config", "", "JSON string with n8n settings (overrides .env, overridden by env vars)")
	flag.Parse()

	if *configStr != "" {
		var cfg configJSON
		if err := json.Unmarshal([]byte(*configStr), &cfg); err != nil {
			log.Fatalf("Failed to parse --config JSON: %v", err)
		}
		setIfNotEnv := func(key, val string) {
			if val != "" && os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
		setIfNotEnv("N8N_API_URL", cfg.N8N_API_URL)
		setIfNotEnv("N8N_API_KEY", cfg.N8N_API_KEY)
	}

	log.Println("Starting Oido n8n MCP Server v1.0.0...")
	RunMCPServer()
}
