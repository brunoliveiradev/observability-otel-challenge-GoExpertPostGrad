package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoExpertPostGrad/observability-otel-challenge-GoExpertPostGrad/service-b-weather-api/internal/model"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
	"net/url"
)

const apiKey = "e1fece5bce574041a9f130048241703"

func GetWeather(location string, ctx context.Context) (*model.WeatherResponse, error) {
	_, span := otel.Tracer("service-b").Start(ctx, "get-weather")
	defer span.End()

	reqStr := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, url.QueryEscape(location))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqStr, nil)
	if err != nil {
		log.Printf("failed to create request to WeatherAPI: %v", err.Error())
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("failed to make request to WeatherAPI: %v", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("failed to get weather from external API, status code: ", resp.StatusCode)
		return nil, fmt.Errorf("failed to get weather")
	}

	var weather model.WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}
