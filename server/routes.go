package server

import (
	"net/http"
	"strings"
	"time"
	"url-shortner/repository"
	"url-shortner/utils"
)

// Router struct to store route handlers
type Router struct {
	mux *http.ServeMux
}

// NewRouter initializes a new router
func NewRouter() *Router {
	return &Router{mux: http.NewServeMux()}
}

// Handle registers a route with a specific HTTP method
func (r *Router) Handle(method string, path string, handler http.HandlerFunc) {
	r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, req)
	})
}

// ServeHTTP makes Router implement http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// DefineRoutes sets up all routes
func DefineRoutes(s *Server) http.Handler {
	router := NewRouter()

	// Define routes with methods
	router.Handle(http.MethodGet, "/s/{id}", handlerGetLongUrl(s.UrlRepository))
	router.Handle(http.MethodPost, "/shortUrlCreate", handlerShortUrlCreate(s.UrlRepository))

	// Apply middlewares
	return LoggingMiddleware(router)
}

// Handler for retrieving the long URL from a short URL ID
func handlerGetLongUrl(urlRepo *repository.URLRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract path param {id}
		// Extract `{id}` manually from URL path
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 || parts[1] != "s" {
			utils.JSONError(w, http.StatusBadRequest, "Invalid URL format")
			return
		}

		idStr := parts[2] // Extracts `{id}` from "/s/{id}"

		if idStr == "" {
			utils.JSONError(w, http.StatusBadRequest, "Missing short URL ID")
			return
		}

		if idStr == "" {
			utils.JSONError(w, http.StatusBadRequest, "Missing short URL ID")
			return
		}

		// Convert base62 to numeric ID
		longID := utils.Base62Decode(idStr)

		// Fetch long URL from database
		longURL, err := urlRepo.AccessLongURL(r.Context(), longID, time.Now())
		if err != nil {
			utils.JSONError(w, http.StatusInternalServerError, "Database error")
			return
		}

		if longURL == nil {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		// Redirect to the long URL
		http.Redirect(w, r, *longURL, http.StatusFound)
	}
}

// Handler for creating a short URL
func handlerShortUrlCreate(urlRepo *repository.URLRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Define request struct
		var req struct {
			LongUrl string `json:"long_url"`
		}

		// Decode JSON request body
		if err := utils.DecodeJSONBody(w, r, &req); err != nil {
			return
		}

		// Validate input
		if req.LongUrl == "" {
			utils.JSONError(w, http.StatusBadRequest, "long_url is required")
			return
		}

		// Save URL to DB
		id, err := urlRepo.SaveURL(r.Context(), req.LongUrl)
		if err != nil {
			utils.JSONError(w, http.StatusInternalServerError, "Failed to create short URL")
			return
		}

		// Respond with the generated ID
		utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
			"short_url": "http://localhost:8080/s/" + utils.Base62Encode(id),
		})
	}
}
