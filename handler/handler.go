package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/sessions"

	"github.com/abdealijaroli/leakybucket/parser"
	"github.com/abdealijaroli/leakybucket/util"
	"github.com/abdealijaroli/leakybucket/db"
)

var database *db.DB

// SetDatabase sets the database connection for the handler package.
func SetDatabase(db *db.DB) {
	database = db
}

var store = sessions.NewCookieStore([]byte("SESSION_SECRET"))

func LinkHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	link := r.PostFormValue("link")
	if link == "" {
		link = session.Values["link"].(string)
	} else {
		session.Values["link"] = link
		session.Save(r, w)
	}

	if r.Header.Get("HX-Request") == "true" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			return
		}
	}

	re, err := regexp.Compile(`^https?://.{1,}\.[^\s/$.?#].[^\s]*$`)
	if err != nil {
		fmt.Fprintf(w, "failed to compile regex: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "https://" + link
	}

	if re.MatchString(link) {
		err := processLink(link)
		if err != nil {
			fmt.Fprintf(w, "failed to process link: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rawData, err := getLinkRawDataFromDB(link)
		if err != nil {
			fmt.Fprintf(w, "failed to get raw data from database: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
				
		res := util.GenerateInitialChatResponse(database.DB, rawData)
		fmt.Fprint(w, res)
		w.WriteHeader(http.StatusOK)
		return

	} else {
		fmt.Fprintf(w, "invalid or no input provided!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func processLink(l string) error {
	err := parser.ParseURL(database.DB, l)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	return err
}

func getLinkRawDataFromDB(link string) (string, error) {
	var rawData string
	err := database.DB.QueryRow("SELECT raw_data FROM scraped_data WHERE link = $1", link).Scan(&rawData)
	if err != nil {
		return "", fmt.Errorf("failed to get raw data from database: %w", err)
	}
	return rawData, nil
}
