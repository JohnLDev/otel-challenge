package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/johnldev/5-validator/config/metrics"
	"github.com/johnldev/5-validator/internal/http/middlewares"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

type Body struct {
	Cep string `json:"cep" validate:"required,len=8"`
}

const serviceName = "validatorApp"

func main() {

	otelShutdown, err := metrics.InitProvider(serviceName, os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatalf("failed to initialize metrics provider: %v", err)
	}

	defer func() {
		if err := otelShutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	tracer := otel.Tracer(serviceName)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "Validation")

		input := Body{}
		json.NewDecoder(r.Body).Decode(&input)
		requestId := middlewares.GetRequestId(r.Context())
		fmt.Printf("request id:%s\n", requestId)

		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(input)

		if err != nil {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		span.End()

		ctx, span = tracer.Start(ctx, "Temperature api call")
		defer span.End()

		span.SetAttributes(attribute.String("cep", input.Cep))

		request, _ := http.NewRequestWithContext(ctx, http.MethodPost, "http://temperature:8080", bytes.NewBuffer([]byte(fmt.Sprintf(`{"zipcode": "%s"}`, input.Cep))))
		request.Header.Set("X-Request-Id", requestId)
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))

		result, err := http.DefaultClient.Do(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer result.Body.Close()

		response, err := io.ReadAll(result.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(string(response))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(result.StatusCode)
		w.Write(response)

	})
	slog.Info("Server started at :8081")
	http.ListenAndServe(":8081", middlewares.PanicRecovery(middlewares.RequestId(mux)))
}
