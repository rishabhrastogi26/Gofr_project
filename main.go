package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gofr.dev/pkg/gofr"
)

// Book represents a book in the bookstore.
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []Book{
	{ID: 1, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"},
	{ID: 2, Title: "To Kill a Mockingbird", Author: "Harper Lee"},
	// Add more books here...
}

func main() {
	app := gofr.New()

	// Define your routes
	app.GET("/books", getBooks)
	app.GET("/books/{id}", getBookByID)
	app.POST("/books", createBook)
	app.PUT("/books/{id}", updateBook)
	app.DELETE("/books/{id}", deleteBook)
	// Add more routes as needed...

	// Start the server
	port := 8080
	fmt.Printf("Server listening on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), app))
}

// Handler for getting all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, books)
}

// Handler for getting a book by ID
func getBookByID(w http.ResponseWriter, r *http.Request) {
	id := gofr.RouteParam(r, "id")
	for _, book := range books {
		if fmt.Sprintf("%d", book.ID) == id {
			respondJSON(w, http.StatusOK, book)
			return
		}
	}
	respondJSON(w, http.StatusNotFound, gofr.Map{"error": "Book not found"})
}

// Handler for creating a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		respondJSON(w, http.StatusBadRequest, gofr.Map{"error": "Invalid request payload"})
		return
	}
	newBook.ID = len(books) + 1
	books = append(books, newBook)
	respondJSON(w, http.StatusCreated, newBook)
}

// Handler for updating an existing book
func updateBook(w http.ResponseWriter, r *http.Request) {
	id := gofr.RouteParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, gofr.Map{"error": "Invalid book ID"})
		return
	}
	var updatedBook Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		respondJSON(w, http.StatusBadRequest, gofr.Map{"error": "Invalid request payload"})
		return
	}
	for i, book := range books {
		if book.ID == bookID {
			books[i] = updatedBook
			respondJSON(w, http.StatusOK, updatedBook)
			return
		}
	}
	respondJSON(w, http.StatusNotFound, gofr.Map{"error": "Book not found"})
}

// Handler for deleting a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	id := gofr.RouteParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, gofr.Map{"error": "Invalid book ID"})
		return
	}
	for i, book := range books {
		if book.ID == bookID {
			books = append(books[:i], books[i+1:]...)
			respondJSON(w, http.StatusOK, gofr.Map{"message": "Book deleted successfully"})
			return
		}
	}
	respondJSON(w, http.StatusNotFound, gofr.Map{"error": "Book not found"})
}

// Helper function to respond with JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
