package main

import (
	"context"
	"log"
	"os"

	"github.com/johnldev/4-deploy-cloud-run/config/metrics"
	"github.com/johnldev/4-deploy-cloud-run/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	defer ctx.Done()

	godotenv.Load()
	shutdown, err := metrics.InitProvider(os.Getenv("SERVICE_NAME"), os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatalf("failed to initialize metrics provider: %s", err.Error())
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown metrics provider: %s", err.Error())
		}
	}()

	server.StartServer()
}
