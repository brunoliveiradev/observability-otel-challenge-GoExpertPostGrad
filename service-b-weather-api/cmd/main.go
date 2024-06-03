package main

import (
	"fmt"
	"github.com/GoExpertPostGrad/observability-otel-challenge-GoExpertPostGrad/service-b-weather-api/internal/cep"
	"github.com/GoExpertPostGrad/observability-otel-challenge-GoExpertPostGrad/service-b-weather-api/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"log"
	"net/http"
)

func initTracer() {
	exporter, err := zipkin.New("http://zipkin:9411/api/v2/spans")
	if err != nil {
		log.Fatalf("Fail to create Zipkin exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("service-b"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func main() {
	initTracer()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/{cep}", func(r chi.Router) {
		r.Use(cep.CheckCepMiddleware)
		r.Get("/", handler.HandleGetTemperatureByCEP)
	})

	fmt.Println("Server running on port 8081")
	http.ListenAndServe(":8081", r)
}
