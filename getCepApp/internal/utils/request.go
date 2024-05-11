package utils

import (
	"context"
	"errors"
	"io"
	"net/http"
	"slices"
)

func RequestWithContext(ctx context.Context, url string) ([]byte, error) {

	// client := &http.Client{
	// 	Transport: otelhttp.NewTransport(&http.Transport{
	// 		TLSClientConfig: &tls.Config{
	// 			InsecureSkipVerify: true,
	// 		},
	// 	},
	// 		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
	// 		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	// 	),
	// 	Timeout: 5 * time.Second,
	// }
	client := http.DefaultClient
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	result, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	response, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	if !slices.Contains([]int{http.StatusOK, http.StatusAccepted, http.StatusCreated}, result.StatusCode) {
		return nil, errors.New(string(response))
	}

	return response, nil
}
