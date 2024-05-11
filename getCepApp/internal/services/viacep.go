package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/johnldev/4-deploy-cloud-run/config"
	"github.com/johnldev/4-deploy-cloud-run/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ICepService interface {
	GetCep(zipcode string) (string, error)
}
type cepService struct {
	ctx context.Context
}

func (s cepService) GetCep(zipcode string) (string, error) {
	tracer := otel.Tracer(config.Conf.ServiceName)
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	spanInitial := trace.SpanFromContext(ctx)

	ctx, span := tracer.Start(ctx, "Get cep api call")
	span.SetAttributes(attribute.String("zipcode", zipcode))
	defer span.End()

	var city string

	resultCdn := make(chan []byte)
	defer close(resultCdn)

	go func() {
		defer utils.PanicRecovery()
		var cepForCdn string = zipcode[:5] + "-" + zipcode[5:]
		cdnUrl := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cepForCdn)
		response, _ := utils.RequestWithContext(ctx, cdnUrl)
		if ctx.Err() == nil {
			resultCdn <- response
		}
	}()

	resultViaCep := make(chan []byte)
	defer close(resultViaCep)

	go func() {
		defer utils.PanicRecovery()
		viaCepUrl := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", strings.Replace(zipcode, "-", "", 1))
		response, _ := utils.RequestWithContext(ctx, viaCepUrl)
		if ctx.Err() == nil {
			resultViaCep <- response
		}
	}()

	for i := 0; i < 2; i++ {
		if ctx.Err() != nil {
			break
		}

		select {
		case result := <-resultCdn:
			fmt.Println("Response from cdn")
			response := struct {
				City string `json:"city"`
			}{}
			json.Unmarshal(result, &response)
			city = response.City
		case result := <-resultViaCep:
			fmt.Println("Response from viacep")
			response := struct {
				City string `json:"localidade"`
			}{}
			json.Unmarshal(result, &response)
			city = response.City
		case <-ctx.Done():
			fmt.Println("Timeout on request")
			fmt.Println(ctx.Err())
		}

		if city != "" {
			cancel()
		}
	}

	if city == "" {
		return "", fmt.Errorf("can not find zipcode")
	}

	spanInitial.SetAttributes(attribute.String("city", city))
	span.SetAttributes(attribute.String("response city", city))

	return city, nil
}

func NewCepService(ctx context.Context) cepService {
	return cepService{ctx: ctx}
}
