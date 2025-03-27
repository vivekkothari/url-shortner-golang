package server

import (
	"log"
	"net/http"
	"url-shortner/db"
	"url-shortner/repository"
)

// Server struct to encapsulate dependencies
type Server struct {
	Addr          string
	db            *db.Database
	UrlRepository *repository.URLRepository
}

// NewServer initializes the server with DB
func NewServer(addr string) (*Server, error) {
	// Initialize DB connection pool
	database, err := db.NewDatabase()
	if err != nil {
		return nil, err
	}

	urlRepository := repository.NewURLRepository(database)

	return &Server{
		Addr:          addr,
		db:            database,
		UrlRepository: urlRepository,
	}, nil
}

// Start runs the HTTP server
func (s *Server) Start() {
	defer s.db.Close() // Ensure DB cleanup on shutdown

	router := DefineRoutes(s) // Pass DB instance to routes

	log.Printf("Starting server on %s...", s.Addr)
	if err := http.ListenAndServe(s.Addr, router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
