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

func GetWeather(location string, ctx context.Context) (*model.WeatherResponse, error) {
	cityEncoded := url.QueryEscape(location)
	return getWeather(cityEncoded, ctx)
}

func NewTemperatureResponse(tempC float64) model.TemperatureResponse {
	return model.TemperatureResponse{
		TempC: tempC,
		TempF: celsiusToFahrenheit(tempC),
		TempK: celsiusToKelvin(tempC),
	}
}

func celsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

func celsiusToKelvin(c float64) float64 {
	return c + 273.15
}

func getWeather(location string, ctx context.Context) (*model.WeatherResponse, error) {
	_, span := otel.Tracer("service-b").Start(ctx, "get-weather")
	defer span.End()

	apiKey := "e1fece5bce574041a9f130048241703"
	formattedUrl := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, location)

	resp, err := http.Get(formattedUrl)
	if err != nil {
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
