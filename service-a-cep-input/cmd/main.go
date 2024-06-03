package main

import (
	"fmt"
	"github.com/GoExpertPostGrad/observability-otel-challenge-GoExpertPostGrad/service-a-cep-input/internal/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"log"
	"net/http"
)

func main() {
	initTracer()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Post("/", handler.HandleCepRequest)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}

func initTracer() {
	exporter, err := zipkin.New("http://zipkin:9411/api/v2/spans")
	if err != nil {
		log.Fatalf("Fail to create Zipkin exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("service-a"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
