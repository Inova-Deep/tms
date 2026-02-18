package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tms/internal/config"
	"tms/internal/db"
	"tms/internal/domain"
	httpHandler "tms/internal/http"
	"tms/internal/logic"
	"tms/seeds"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()

	database, err := db.Init(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Printf("Connected to database at %s", cfg.DBPath)

	if err := seeds.Run(database); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	// Initialize repositories (single instances, properly shared)
	employeeRepo := &domain.EmployeeRepository{DB: database}
	profileRepo := &domain.ProfileRepository{DB: database}
	evidenceRepo := &domain.EvidenceRepository{DB: database}
	eventRepo := &domain.EventRepository{DB: database}
	certificationRepo := &domain.CertificationRepository{DB: database}
	courseRepo := &domain.CourseRepository{DB: database}

	// Initialize services
	computationService := &logic.ComputationService{}
	dashboardService := &logic.DashboardService{
		EmployeeRepo:       employeeRepo,
		ProfileRepo:        profileRepo,
		EvidenceRepo:       evidenceRepo,
		ComputationService: computationService,
	}

	// Create HTTP server
	srv := &httpHandler.Server{
		EmployeeRepo:       employeeRepo,
		ProfileRepo:        profileRepo,
		EvidenceRepo:       evidenceRepo,
		EventRepo:          eventRepo,
		CertificationRepo:  certificationRepo,
		CourseRepo:         courseRepo,
		ComputationService: computationService,
		DashboardService:   dashboardService,
	}

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS configuration (demo: allow localhost origins)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Mount("/api", srv.Routes())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("TMS Backend Running"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Graceful shutdown
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
