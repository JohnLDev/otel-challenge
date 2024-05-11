package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/johnldev/4-deploy-cloud-run/config"
	"github.com/johnldev/4-deploy-cloud-run/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	weatherApiUrl = "http://api.weatherapi.com/v1/current.json"
)

type ITemperatureResponse struct {
	Current struct {
		Celcius    float64 `json:"temp_c"`
		Fahrenheit float64 `json:"temp_f"`
	} `json:"current"`
}

type IWeatherService interface {
	GetTemperatureByCity(city string) (*ITemperatureResponse, error)
}

type WeatherService struct {
	Ctx context.Context
}

func (s WeatherService) GetTemperatureByCity(city string) (*ITemperatureResponse, error) {
	tracer := otel.Tracer(config.Conf.ServiceName)

	spanInitial := trace.SpanFromContext(s.Ctx)

	ctx, span := tracer.Start(s.Ctx, "Weather api call")
	defer span.End()

	span.SetAttributes(attribute.String("city", city))

	token := os.Getenv("WEATHER_API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("weather api token not found")
	}

	result := ITemperatureResponse{}
	cityInput := utils.NormalizeAccents(strings.ToLower(strings.ReplaceAll(city, " ", "-")))

	response, err := utils.RequestWithContext(ctx, fmt.Sprintf("%s?q=%s&key=%s", weatherApiUrl, cityInput, token))
	if err != nil {
		return nil, err
	}

	//todo validate error

	json.Unmarshal(response, &result)
	span.SetAttributes(attribute.Float64("Celcius result", result.Current.Celcius))
	spanInitial.SetAttributes(attribute.Float64("Celcius result", result.Current.Celcius))
	span.SetAttributes(attribute.Float64("Fahrenheit result", result.Current.Fahrenheit))
	spanInitial.SetAttributes(attribute.Float64("Fahrenheit result", result.Current.Fahrenheit))
	return &result, nil
}

func NewWeatherService(ctx context.Context) WeatherService {
	return WeatherService{
		Ctx: ctx,
	}
}
