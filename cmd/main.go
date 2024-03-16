package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/abdealijaroli/leakybucket/handler"
)

func main() {
	fmt.Println(`
+---_-----------_--------,-__------------_-----------+
| \_|_)         | |      /|/  \          | |         |
|   |    _  __, | |       | __/       __ | |  __|_   |
|  _|   |/ /  | |/_) |   ||   \|   | /   |/_)|/ |    |
| (/\___/__|_/|_/ \_/ \_/|/(__/ \_/|_|___/ \_/__/_/  |
|                       /|                           |
|                       \|                           |
+-----------------------------------------------------+
					   `)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	fs := http.FileServer(http.Dir("./view"))
	http.Handle("/", fs)

	http.HandleFunc("/link", handler.LinkHandler)

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
