package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/abdealijaroli/leakybucket/parser"
	"github.com/abdealijaroli/leakybucket/util"
)

type Handler struct {
	DB *sql.DB
}

func (h *Handler) LinkHandler(w http.ResponseWriter, r *http.Request) {
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

	link := r.PostFormValue("link")

	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "https://" + link
	}

	if re.MatchString(link) {
		err := h.processLink(link)
		if err != nil {
			fmt.Fprintf(w, "failed to process link: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rawData, err := h.getLinkRawDataFromDB(link)
		if err != nil {
			fmt.Fprintf(w, "failed to get raw data from database: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
				
		res := util.GenerateInitialChatResponse(h.DB, rawData)
		fmt.Fprint(w, res)
		w.WriteHeader(http.StatusOK)
		return

	} else {
		fmt.Fprintf(w, "invalid or no input provided!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) processLink(l string) error {
	err := parser.ParseURL(h.DB, l)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	return err
}

func (h *Handler) getLinkRawDataFromDB(link string) (string, error) {
	var rawData string
	err := h.DB.QueryRow("SELECT raw_data FROM scraped_data WHERE link = $1", link).Scan(&rawData)
	if err != nil {
		return "", fmt.Errorf("failed to get raw data from database: %w", err)
	}
	return rawData, nil
}
