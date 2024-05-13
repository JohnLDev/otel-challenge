package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/johnldev/4-deploy-cloud-run/config"
	"github.com/johnldev/4-deploy-cloud-run/internal/services"
	usecases "github.com/johnldev/4-deploy-cloud-run/internal/useCases"
	"github.com/johnldev/4-deploy-cloud-run/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Input struct {
	Zipcode string `json:"zipcode" validate:"required,len=8,number"`
}

func StartServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Use(otelmux.Middleware(config.Conf.ServiceName))
	r.Use(middleware.Recoverer)
	r.Use(utils.RequestIDMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer(config.Conf.ServiceName)

		ctx, span := tracer.Start(r.Context(), "GetCepApp received request")
		defer span.End()

		validate := validator.New(validator.WithRequiredStructEnabled())
		input := Input{}
		// Decode the request body into the input struct
		json.NewDecoder(r.Body).Decode(&input)

		validateErr := validate.Struct(input)
		if validateErr != nil {
			http.Error(w, "invalid zipcode", http.StatusBadRequest)
			return
		}
		span.SetAttributes(attribute.String("zipcode", input.Zipcode))

		response, err := usecases.NewGetTemperatureUseCase(services.NewCepService(ctx), services.NewWeatherService(ctx)).Execute(ctx, input.Zipcode)
		if err != nil {
			if err.Error() == "can not find zipcode" {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"city": %s,"temp_C": %.2f, "temp_F": %.2f, "temp_K": %.2f}`, response.City, response.Celcius, response.Fahrenheit, response.Kelvin)))
	})
	slog.Info("Server started at :8080")
	http.ListenAndServe(":8080", r)
}
