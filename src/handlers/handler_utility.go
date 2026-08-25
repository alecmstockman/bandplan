package handlers

import (
	"log"
	"net/http"
)

func (h Handler) HandlerHealth(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerHealth")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("ok")); err != nil {
		log.Println("Unable to write health-check response:", err)
	}
}
