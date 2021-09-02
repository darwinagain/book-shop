package postgres

import (
	"book-shop/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"github.com/pquerna/ffjson/ffjson"
)

var (
	ErrNotFound      = errors.New("book not found")
	ErrDuplicateBook = errors.New("duplicate book found")
)

func GetAllBooks() ([]*models.Book, error) {
	db := OpenConnection()

	var books []*models.Book

	query := ("SELECT id, title, author_first_name, author_last_name, genre, quantity FROM books")

	rows, err := db.Query(query)
	if err != nil {
		return []*models.Book{}, err
	}

	defer rows.Close()

	for rows.Next() {
		book := new(models.Book)
		if err := rows.Scan(&book.ID, &book.Title, &book.AuthorFirstName, &book.AuthorLastName, &book.Genre, &book.Quantity); err != nil {
			return []*models.Book{}, err
		}
		books = append(books, book)
	}

	defer db.Close()

	return books, nil
}

func GetBookByID(bookID int) (models.Book, error) {
	db := OpenConnection()

	var bookData models.Book

	query := "SELECT id, title, author_first_name, author_last_name, genre, quantity FROM books WHERE id = $1"

	err := db.QueryRow(query, bookID).Scan(&bookData.ID, &bookData.Title, &bookData.AuthorFirstName, &bookData.AuthorLastName, &bookData.Genre, &bookData.Quantity)
	if err != nil {
		return models.Book{}, err
	}

	defer db.Close()

	return bookData, nil
}

func CreateBook(bookData models.Book) ([]byte, error) {
	db := OpenConnection()

	existingBook, err := bookExists(db, bookData.Title, bookData.AuthorFirstName, bookData.AuthorLastName)
	if err != nil {
		return []byte{}, err
	}

	if existingBook {
		return []byte{}, ErrDuplicateBook
	}

	err = db.QueryRow("SELECT id FROM books ORDER BY id DESC LIMIT 1").Scan(&bookData.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			return []byte{}, err
		}
		bookData.ID = 0
	}

	bookData.ID += 1

	query := "INSERT INTO books (title, author_first_name, author_last_name, genre, quantity) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.Exec(query, bookData.Title, bookData.AuthorFirstName, bookData.AuthorLastName, bookData.Genre, bookData.Quantity)
	if err != nil {
		return []byte{}, err
	}

	bookBytes, _ := ffjson.Marshal(bookData)

	return bookBytes, nil
}

func UpdateBook(bookID int, bookData *models.Book) ([]byte, error) {
	var book models.Book

	db := OpenConnection()

	// TODO: Add check to see if book with given id exists

	query := sq.Update("books").Where(sq.Eq{"id": bookID})

	if bookData.Title != "" {
		query = query.Set("title", bookData.Title)
	}
	if bookData.AuthorFirstName != "" {
		query = query.Set("author_first_name", bookData.AuthorFirstName)
	}
	if bookData.AuthorLastName != "" {
		query = query.Set("author_last_name", &bookData.AuthorLastName)
	}
	if bookData.Genre != "" {
		query = query.Set("genre", &bookData.Genre)
	}
	if bookData.Quantity != nil {
		query = query.Set("quantity", &bookData.Quantity)
	}

	query = query.Suffix("RETURNING id, title, author_first_name, author_last_name, genre, quantity")

	newQuery, args, err := query.ToSql()

	cleanQuery := replaceSQLPlaceholder(newQuery)

	if err != nil {
		return []byte{}, errors.New("error building books query")
	}

	row := db.QueryRow(cleanQuery, args...)

	err = row.Scan(&book.ID, &book.Title, &book.AuthorFirstName, &book.AuthorLastName, &book.Genre, &book.Quantity)
	if err != nil {
		return []byte{}, err
	}

	newBookBytes, _ := ffjson.Marshal(book)

	return newBookBytes, nil
}

func DeleteBook(bookID int) error {
	db := OpenConnection()

	query := "DELETE FROM books WHERE id = $1"

	_, err := db.Exec(query, bookID)
	if err != nil {
		return err
	}

	return nil
}

func bookExists(db *sql.DB, title string, authorFirstName string, authorLastName string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM books WHERE title = $1 AND author_first_name = $2 AND author_last_name = $3)"

	row := db.QueryRow(query,
		title, authorFirstName, authorLastName)
	if err := row.Scan(&exists); err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
			return false, err
		}

		return exists, nil
	}

	return exists, nil
}

func replaceSQLPlaceholder(sql string) string {
	for nParam := 1; strings.Contains(sql, "?"); nParam++ {
		sql = strings.Replace(sql, "?", fmt.Sprintf("$%d", nParam), 1)
	}
	return sql
}
