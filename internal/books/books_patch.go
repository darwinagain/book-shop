package books

import (
	"book-shop/internal/models"
	"book-shop/internal/postgres"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func updateBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Endpoint Hit: updateBook")

		vars := mux.Vars(r)
		bookID := vars["id"]
		bookIDInt, _ := strconv.Atoi(bookID)

		// ctx := r.Context()

		var reqBody models.Book

		// decode json body
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// store in postgres
		output, err := postgres.UpdateBook(bookIDInt, &reqBody)
		if err != nil {
			log.Printf("Error storing to database: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(output); err != nil {
			log.Printf("Error writing payload: %v", err)
		}
	}
}
