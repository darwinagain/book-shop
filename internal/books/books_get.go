package books

import (
	"book-shop/internal/models"
	"book-shop/internal/postgres"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pquerna/ffjson/ffjson"
)

func returnAllBooks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Endpoint Hit: returnAllBooks")

		books, err := postgres.GetAllBooks()
		if err != nil {
			log.Printf("Error getting books: %v", err)
			return
		}

		output, _ := ffjson.Marshal(books)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(output); err != nil {
			log.Printf("Error writing payload: %v", err)
		}
	}
}

func returnSingleBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Endpoint Hit: returnSingleBook")

		vars := mux.Vars(r)
		bookID := vars["id"]
		bookIDInt, _ := strconv.Atoi(bookID)

		book, err := postgres.GetBookByID(bookIDInt)
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

		output, _ := ffjson.Marshal(book)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(output); err != nil {
			log.Printf("Error writing payload: %v", err)
		}
	}
}

func returnQueryResults() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Endpoint Hit: returnQueryResults")

		bookID, _ := strconv.Atoi(r.URL.Query().Get("id"))
		title := r.URL.Query().Get("title")
		author := r.URL.Query().Get("author")
		genre := r.URL.Query().Get("genre")

		book := models.Book2{
			ID:     bookID,
			Title:  title,
			Author: author,
			Genre:  genre,
		}

		books, err := postgres.GetQueryResults(book)
		if err != nil {
			log.Printf("Error getting books: %v", err)
			return
		}

		output, _ := ffjson.Marshal(books)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(output); err != nil {
			log.Printf("Error writing payload: %v", err)
		}
	}
}
