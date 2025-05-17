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
	Router        *http.Handler
}

// NewServer initializes the server with DB
func NewServer(addr string) (*Server, error) {
	// Initialize DB connection pool
	database, err := db.NewDatabase()
	if err != nil {
		return nil, err
	}

	urlRepository := repository.NewURLRepository(database)

	server := Server{
		Addr:          addr,
		db:            database,
		UrlRepository: urlRepository,
	}
	router := DefineRoutes(&server) // Pass DB instance to routes
	server.Router = &router
	return &server, nil
}

// Start runs the HTTP server
func (s *Server) Start() {
	defer s.db.Close() // Ensure DB cleanup on shutdown

	log.Printf("Starting server on %s...", s.Addr)
	if err := http.ListenAndServe(s.Addr, *s.Router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
