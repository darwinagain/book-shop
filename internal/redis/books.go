package redis

import (
	"book-shop/internal/models"
	"bytes"
	"sort"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/pquerna/ffjson/ffjson"
)

// Redis keys
const (
	// LastBookID holds id of the last book that was created
	LastBookID = "last_book_id"
	// BookList holds ids of all books that have been created
	BookList = "book_list"
)

// Declare a pool variable to hold the pool of Redis connections
var pool = newPool()

func GetAllBooks() ([]models.Book, error) {
	conn := pool.Get()
	defer conn.Close()

	// Get the list of all book ids
	bookIDs, err := redis.Ints(conn.Do("SMEMBERS", BookList))
	if err != nil {
		return []models.Book{}, err
	}

	keys := BookKeysByID(bookIDs)

	// Get all book data using ids found above
	bookInterfaces, err := redis.ByteSlices(conn.Do("MGET", keys...))
	if err != nil {
		return []models.Book{}, err
	}

	// Unmarshal books
	books, err := groupUnmarshal(bookInterfaces)
	if err != nil {
		return []models.Book{}, err
	}

	return books, nil
}

func GetBookByID(bookID int) (models.Book, error) {
	conn := pool.Get()
	defer conn.Close()

	// Get book using id supplied during api call
	bookBytes, err := redis.Bytes(conn.Do("GET", BookDataKey(bookID)))
	if err != nil {
		return models.Book{}, err
	}

	var bookData models.Book

	// Unmarshal book data
	err = ffjson.Unmarshal(bookBytes, &bookData)
	if err != nil {
		return models.Book{}, err
	}

	return bookData, nil
}

func CreateBook(bookData models.Book) ([]byte, error) {
	conn := pool.Get()
	defer conn.Close()

	// Generate a new book id by incrementing from the last book id
	nextBookID, err := redis.Int(conn.Do("INCR", LastBookID))
	if err != nil {
		return []byte{}, err
	}

	// Set the book id
	bookData.ID = nextBookID

	// Marshal the book data
	newBookBytes, _ := ffjson.Marshal(bookData)

	// Start a multi request to redis
	conn.Send("MULTI")

	conn.Send("SET", BookDataKey(nextBookID), newBookBytes)

	// Add book id to set of all books
	conn.Send("SADD", BookList, nextBookID)
	_, err = conn.Do("EXEC")
	if err != nil {
		return []byte{}, err
	}

	return newBookBytes, nil
}

func BookKeysByID(ids []int) []interface{} {
	keys := make([]interface{}, len(ids))
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	for i, id := range ids {
		keys[i] = BookDataKey(id)
	}
	return keys
}

func BookDataKey(id int) string {
	return "book:" + strconv.Itoa(id)
}

func groupUnmarshal(byteSlices [][]byte) ([]models.Book, error) {
	cleanedResp := make([][]byte, 0, len(byteSlices))

	for _, byteSlice := range byteSlices {
		if byteSlice != nil {
			cleanedResp = append(cleanedResp, byteSlice)
		}
	}

	joinedOutput := bytes.Join(cleanedResp, []byte(","))
	joinedJSONBytes := bytes.Join([][]byte{
		[]byte("["),
		joinedOutput,
		[]byte("]"),
	}, []byte(""))

	var out []models.Book
	err := ffjson.Unmarshal(joinedJSONBytes, &out)
	if err != nil {
		return []models.Book{}, err
	}

	return out, nil
}
