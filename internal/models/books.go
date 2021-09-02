package models

type Book struct {
	ID              int    `json:"id" db:"id"`
	Title           string `json:"title" db:"title" validate:"required,ascii,max=128"`
	AuthorFirstName string `json:"author_first_name" db:"author_first_name" validate:"required,ascii,max=128"`
	AuthorLastName  string `json:"author_last_name"  db:"author_last_name" validate:"required,ascii,max=128"`
	Genre           string `json:"genre"  db:"genre" validate:"required,ascii,max=128"`
	Quantity        *int   `json:"quantity" db:"quantity" validate:"required,numeric"`
}
