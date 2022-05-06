package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"book-shop/graph/generated"
	"book-shop/graph/model"
	"book-shop/internal/models"
	"book-shop/internal/postgres"
	"context"
	"log"
	"strconv"

	"github.com/pquerna/ffjson/ffjson"
)

func (r *mutationResolver) CreateBook(ctx context.Context, input model.NewBook) (*model.Book, error) {

	book := models.Book{
		Title:           input.Title,
		AuthorFirstName: input.AuthorFirstName,
		AuthorLastName:  input.AuthorLastName,
		Genre:           input.Genre,
	}

	output, err := postgres.CreateBook(book)

	if err != nil {
		if err == postgres.ErrDuplicateBook {
			msg := "book already exists"
			log.Println(msg)
		}
		log.Printf("Error storing to database: %v", err)
		return &model.Book{}, err
	}

	var newBook models.Book

	err = ffjson.Unmarshal(output, &newBook)
	if err != nil {
		return &model.Book{}, err
	}

	finalBook := model.Book{
		ID:              strconv.Itoa(newBook.ID),
		Title:           newBook.Title,
		AuthorFirstName: newBook.AuthorFirstName,
		AuthorLastName:  newBook.AuthorLastName,
		Genre:           newBook.Genre,
	}

	return &finalBook, nil
}

func (r *queryResolver) Books(ctx context.Context) ([]*model.Book, error) {
	var books []*model.Book

	// get all books using existing postgres function
	_books, err := postgres.GetAllBooks()
	if err != nil {
		log.Printf("Error getting books: %v", err)
		return books, err
	}

	// convert each book from Go Book type to GQL Book type and append to books
	for _, v := range _books {
		book := &model.Book{
			ID:              strconv.Itoa(v.ID),
			Title:           v.Title,
			AuthorFirstName: v.AuthorFirstName,
			AuthorLastName:  v.AuthorLastName,
			Genre:           v.Genre,
		}
		books = append(books, book)
	}

	return books, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
