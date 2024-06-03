package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoExpertPostGrad/observability-otel-challenge-GoExpertPostGrad/service-b-weather-api/internal/model"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
)

func GetAddressFromViaCEP(cep string, ctx context.Context) (*model.AddressResponse, error) {
	_, span := otel.Tracer("service-b").Start(ctx, "get-cep-location")
	defer span.End()

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("failed to get location from ViaCEP, status code: ", resp.StatusCode)
		return nil, fmt.Errorf("failed to get location")
	}

	var address model.AddressResponse
	err = json.NewDecoder(resp.Body).Decode(&address)
	if address.Erro {
		return nil, fmt.Errorf("zipcode not found")
	}

	if err != nil {
		return nil, err
	}

	if address.Localidade == "" {
		return nil, fmt.Errorf("can not find zipcode")
	}

	return &address, nil
}
