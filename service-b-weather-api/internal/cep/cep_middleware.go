package cep

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func CheckCepMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cep := chi.URLParam(r, "cep")
		if cep == "" || !isValidZipcode(cep) {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isValidZipcode(zipcode string) bool {
	if len(zipcode) != 8 {
		return false
	}
	for _, char := range zipcode {
		if _, err := strconv.Atoi(string(char)); err != nil {
			return false
		}
	}
	return true
}
