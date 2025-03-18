package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DonBigBon/parser-backend/config"
	handlers "github.com/DonBigBon/parser-backend/internal/api"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	router.HandleFunc("/upload", handlers.UploadHandler).Methods("POST")
	router.HandleFunc("/download", handlers.DownloadHandler).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
