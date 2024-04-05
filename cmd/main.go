package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/abdealijaroli/leakybucket/db"
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
	
	//load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file: ", err)
	}

	//connect db
	database, err := db.NewDB()
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
	defer database.Close()

	handler.SetDatabase(database)

	// seed db (one time operation)
	// err = db.SeedDB(database)
	// if err != nil {
	// 	log.Fatal("error while seeding data to database: ", err)
	// }

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

	log.Println("server is running on port", port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("error starting server: ", err)
	}
}
