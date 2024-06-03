package handler

import (
	"encoding/json"
	"github.com/GoExpertPostGrad/observability-otel-challenge-GoExpertPostGrad/service-b-weather-api/internal/service"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
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
		http.Error(w, "can not find weather", http.StatusNotFound)
		return
	}

	temperatureResponse := service.NewTemperatureResponse(weather.Current.TempC)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(temperatureResponse)
}
