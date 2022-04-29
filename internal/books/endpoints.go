package books

import (
	"book-shop/internal/models"
	"context"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gorilla/mux"
)

func Run(ctx context.Context) {
	log.Println("starting server")

	handleRequests()
}

func handleRequests() {

	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/health/live", getLiveness()).Methods("GET")
	myRouter.HandleFunc("/books/all", returnAllBooks()).Methods("GET")
	myRouter.HandleFunc("/books", returnQueryResults()).Methods("GET")
	myRouter.HandleFunc("/book/{id}", returnSingleBook()).Methods("GET")
	myRouter.HandleFunc("/book/{id}", updateBook()).Methods("PATCH")
	myRouter.HandleFunc("/book/{id}", deleteBook()).Methods("DELETE")
	myRouter.HandleFunc("/book", createNewBook()).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func validateRequest(ctx context.Context, req models.Book) map[string]string {

	errorResponse := make(map[string]string, 30)
	errorField := make([]string, 0, 10)

	translator := en.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := v.Struct(req)
	if err == nil {
		return errorResponse
	}

	for _, e := range err.(validator.ValidationErrors) {
		errorField = append(errorField, e.Translate(trans))
	}

	if len(errorField) > 0 {
		errorResponse["message"] = strings.Join(errorField, ", ")
	}

	return errorResponse
}
