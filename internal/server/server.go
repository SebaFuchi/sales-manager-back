package server

import (
	"log"
	"net/http"
	"sales-manager-back/internal/server/routes"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"
	"sales-manager-back/pkg/useCases/Helpers/firebaseHelper"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server struct holds the HTTP server instance
type Server struct {
	server *http.Server
}

// New initializes and configures the HTTP server with security middleware
func New(port string) (*Server, error) {
	r := chi.NewRouter()

	// Security & Performance Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration - restrict to your frontend domain in production
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://localhost:*", "http://127.0.0.1:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Tenant-ID", "X-User-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))

	// Mount API routes
	r.Mount("/api", routes.New())

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	serv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	server := Server{server: serv}
	return &server, nil
}

func (serv *Server) Start() {
	// Initialize database connection and run migrations
	dbConn := databaseHelper.InitDB()
	databaseHelper.Db = dbConn
	log.Println("✓ Database connected and migrations executed")

	// Initialize Firebase Admin SDK
	firebaseHelper.InitFirebase()

	log.Printf("✓ Sales Manager API server starting on %s", serv.server.Addr)
	log.Printf("✓ Health check: http://localhost%s/health", serv.server.Addr)
	log.Printf("✓ API base: http://localhost%s/api/sales-manager", serv.server.Addr)

	if err := serv.server.ListenAndServe(); err != nil {
		log.Fatalf("✗ Server failed: %v", err)
	}
}
