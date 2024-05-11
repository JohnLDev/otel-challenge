package usecases

import (
	"context"
	"fmt"

	"github.com/johnldev/4-deploy-cloud-run/internal/services"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type GetTemperatureUseCase struct {
	CepService     services.ICepService
	WeatherService services.IWeatherService
}

type GetTemperaturerResponse struct {
	Celcius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func (u *GetTemperatureUseCase) Execute(ctx context.Context, zipcode string) (*GetTemperaturerResponse, error) {
	// call via zipcode api
	city, err := u.CepService.GetCep(zipcode)
	if err != nil {
		return nil, err
	}
	fmt.Println(city)

	tempResponse, err := u.WeatherService.GetTemperatureByCity(city)
	if err != nil {
		return nil, err
	}

	response := GetTemperaturerResponse{
		Celcius:    tempResponse.Current.Celcius,
		Fahrenheit: tempResponse.Current.Fahrenheit,
	}

	response.Kelvin = tempResponse.Current.Celcius + 273.15
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Float64("Kelvin response", response.Kelvin))
	return &response, nil

}

func NewGetTemperatureUseCase(cepService services.ICepService, weatherService services.IWeatherService) *GetTemperatureUseCase {
	return &GetTemperatureUseCase{
		CepService:     cepService,
		WeatherService: weatherService,
	}
}
