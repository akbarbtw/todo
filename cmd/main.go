package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"todo-list/internal/handler"
	"todo-list/internal/repository"
	"todo-list/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db := initDB()
	defer db.Close()

	repo := repository.NewTaskRepository(db)
	svc := service.NewTaskService(repo)
	taskHandler := handler.NewTaskHandler(svc)

	router := mux.NewRouter()

	// CORS Middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	taskHandler.RegisterRoutes(router)

	// Serve Static Frontend Web UI
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

func initDB() *sql.DB {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "todo")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	var db *sql.DB
	var err error

	// Retry connection for up to 30 seconds
	for i := 0; i < 30; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Failed to open database: %v, retrying...", err)
			time.Sleep(time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Failed to ping database: %v, retrying...", err)
			time.Sleep(time.Second)
			continue
		}

		log.Println("Database connected successfully")
		break
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after retries: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
