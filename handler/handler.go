package handler

import (
	"fmt"
	"net/http"

	"github.com/abdealijaroli/leakybucket/util"
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
		response := util.GenerateLinkResponse(link)
		fmt.Fprintf(w, "Response: %s", response)
	} else {
		fmt.Fprintf(w, "No input provided!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
