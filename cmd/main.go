package main

import (
	"book-shop/internal/books"
	"context"
)

func main() {
	books.Run(context.Background())
}
