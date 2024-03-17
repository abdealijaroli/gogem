package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/abdealijaroli/leakybucket/handler"
	"github.com/abdealijaroli/leakybucket/db"
)

func main() {
	//load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	//connect db
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// serve static files
	fs := http.FileServer(http.Dir("./view"))
	http.Handle("/", fs)

	// handle api routes
	http.HandleFunc("/link", handler.LinkHandler)

	// start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server is running on port", port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
