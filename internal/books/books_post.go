package books

import (
	"book-shop/internal/models"
	"book-shop/internal/postgres"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pquerna/ffjson/ffjson"
)

func createNewBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Endpoint Hit: createNewBook")

		ctx := r.Context()

		var reqBody models.Book

		fmt.Println(r.Body)

		// decode json body
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// validate request
		errs := validateRequest(ctx, reqBody)
		if len(errs) > 0 {
			log.Printf("Error parsing request: %v", errs)

			outMap := map[string]interface{}{
				"error": errs,
			}

			output, _ := ffjson.Marshal(outMap)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write(output); err != nil {
				log.Printf("Error writing payload: %v", err)
			}
			return
		}

		// store in postgres
		output, err := postgres.CreateBook(reqBody)
		if err != nil {
			if err == postgres.ErrDuplicateBook {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				msg := "book already exists"
				if _, err := w.Write([]byte(msg)); err != nil {
					log.Printf("Error writing payload: %v", err)
				}
			}
			log.Printf("Error storing to database: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(output); err != nil {
			log.Printf("Error writing payload: %v", err)
		}
	}
}
