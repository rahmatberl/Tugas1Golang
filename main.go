package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type Book struct {
	BookID        string `json:"BookID"`
	BookName      string `json:"BookName"`
	BookCategory  string `json:"BookCategory"`
	BookYear      string `json:"BookYear"`
	BookAuthor    string `json:"BookAuthor"`
	BookPublisher string `json:"BookPublisher"`
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []Book

	sql := `SELECT
				BookID,
				IFNULL(BookName,''),
				IFNULL(BookCategory,'') BookCategory,
				IFNULL(BookYear,'') BookYear,
				IFNULL(BookAuthor,'') BookAuthor,
				IFNULL(BookPublisher,'') BookPublisher
			FROM books`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var book Book
		err := result.Scan(&book.BookID, &book.BookName, &book.BookCategory, &book.BookYear, &book.BookAuthor, &book.BookPublisher)

		if err != nil {
			panic(err.Error())
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		bookID := r.FormValue("BookID")
		bookName := r.FormValue("BookName")
		bookCategory := r.FormValue("BookCategory")
		bookYear := r.FormValue("BookYear")
		bookAuthor := r.FormValue("BookAuthor")
		bookPublisher := r.FormValue("BookPublisher")

		stmt, err := db.Prepare("INSERT INTO books(BookID,BookName,BookCategory,BookYear,BOokAuthor,BookPublisher) VALUES (?,?,?,?,?,?)")

		_, err = stmt.Exec(bookID, bookName, bookCategory, bookYear, bookAuthor, bookPublisher)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}
	}
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var books []Book
	params := mux.Vars(r)

	sql := `SELECT
				BookID,
				IFNULL(BookName,''),
				IFNULL(BookCategory,'') BookCategory,
				IFNULL(BookYear,'') BookYear,
				IFNULL(BookAuthor,'') BookAuthor,
				IFNULL(BookPublisher,'') BookPublisher
			FROM books`

	result, err := db.Query(sql, params["id"])

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var book Book

	for result.Next() {
		err := result.Scan(&book.BookID, &book.BookName, &book.BookCategory, &book.BookYear, &book.BookAuthor, &book.BookPublisher)

		if err != nil {
			panic(err.Error())
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		params := mux.Vars(r)

		newBookName := r.FormValue("BookName")
		newBookCategory := r.FormValue("BookCategory")
		newBookYear := r.FormValue("BookYear")
		newBookAuthor := r.FormValue("BookAuthor")
		newBookPublisher := r.FormValue("BookPublisher")

		stmt, err := db.Prepare("UPDATE books set BookName = ?, BookCategory = ?, BookYear = ?, BookAuthor = ?, BookPublisher = ? WHERE BookID = ?")

		_, err = stmt.Exec(newBookName, newBookCategory, newBookYear, newBookAuthor, newBookPublisher, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data Not Found or Request Error")
		}

		fmt.Fprintf(w, "Book with BookID = %s was Updated", params["id"])
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM books WHERE BookID =?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "Delete Failed")
	}

	fmt.Fprintf(w, "Book with BookID = %s was Deleted", params["id"])
}

func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []Book

	bookID := r.FormValue("BookID")
	bookName := r.FormValue("BookName")

	sql := `SELECT
				BookID,
				IFNULL(BookName,''),
				IFNULL(BookCategory,'') BookCategory,
				IFNULL(BookYear,'') BookYear,
				IFNULL(BookAuthor,'') BookAuthor,
				IFNULL(BookPublisher,'') BookPublisher
			FROM books WHERE BookID = ? AND BookName = ?`

	result, err := db.Query(sql, bookID, bookName)

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var book Book

	for result.Next() {
		err := result.Scan(&book.BookID, &book.BookName, &book.BookCategory, &book.BookYear, &book.BookAuthor, &book.BookPublisher)

		if err != nil {
			panic(err.Error())
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func main() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/bookstore")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/books", getBook).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	//New
	r.HandleFunc("/getcustomer", getPost).Methods("POST")

	log.Fatal(http.ListenAndServe(":4321", r))
}
