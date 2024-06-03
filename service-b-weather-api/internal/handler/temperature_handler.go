package handler

import (
	"encoding/json"
	"github.com/GoExpertPostGrad/observability-otel-challenge-GoExpertPostGrad/service-b-weather-api/internal/model"
	"github.com/GoExpertPostGrad/observability-otel-challenge-GoExpertPostGrad/service-b-weather-api/internal/service"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
)

func HandleGetTemperatureByCEP(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	ctx, span := otel.Tracer("service-b").Start(r.Context(), "get-cep-temperature")
	defer span.End()

	address, err := service.GetAddressFromViaCEP(cep, ctx)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	weather, err := service.GetWeather(address.Localidade, ctx)
	if err != nil {
		log.Printf("Error fetching weather: %v", err)
		http.Error(w, "can not find weather", http.StatusNotFound)
		return
	}

	log.Printf("Weather fetched: %+v", weather)

	temperature := model.TemperatureResponse{
		City:  address.Localidade,
		TempC: weather.Current.TempC,
		TempF: celsiusToFahrenheit(weather.Current.TempC),
		TempK: celsiusToKelvin(weather.Current.TempC),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(temperature)
}

func celsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

func celsiusToKelvin(c float64) float64 {
	return c + 273.15
}
