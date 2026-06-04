package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	"bitbucket_mcp/internal/bitbucket"
	"bitbucket_mcp/internal/config"
	"bitbucket_mcp/internal/server"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	client := bitbucket.NewClient(cfg.Workspace, cfg.Email, cfg.Token)
	srv := server.New(client)

	if err := srv.Run(context.Background()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
