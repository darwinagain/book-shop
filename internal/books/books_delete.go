package books

import (
	"book-shop/internal/postgres"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func deleteBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Endpoint Hit: returnSingleBook")

		vars := mux.Vars(r)
		bookID := vars["id"]
		bookIDInt, _ := strconv.Atoi(bookID)

		err := postgres.DeleteBook(bookIDInt)
		if err != nil {
			log.Println(err)
			msg := fmt.Sprintf("book not found: id %s", bookID)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte(msg)); err != nil {
				log.Printf("Error writing payload: %v", err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}
