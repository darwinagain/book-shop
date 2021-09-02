package books

import (
	"log"
	"net/http"
)

func getLiveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Endpoint Hit: getLiveness")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
