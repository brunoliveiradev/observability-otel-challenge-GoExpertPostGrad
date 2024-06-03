package service

import (
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel"
	"io"
	"net/http"
)

type TemperatureResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func GetTemperature(cep string, ctx context.Context) (*TemperatureResponse, int, error) {
	_, span := otel.Tracer("service-a").Start(ctx, "request-service-b")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://goapp-service-b:8081/"+cep, nil)
	if err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	var temperatureResponse TemperatureResponse
	err = json.Unmarshal(body, &temperatureResponse)
	if err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	return &temperatureResponse, http.StatusOK, nil
}
