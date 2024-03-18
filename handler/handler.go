package handler

import (
	"fmt"
	"net/http"

	"github.com/abdealijaroli/leakybucket/parser"
)

func LinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
	}

	link := r.PostFormValue("link")
	if link != "" {
		resp, err := parser.ParseURL(link)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		} else {
			fmt.Fprintf(w, "%s", resp)
		}
	} else {
		fmt.Fprintf(w, "No input provided!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
