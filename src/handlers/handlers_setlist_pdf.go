package handlers

import (
	"log"
	"net/http"
)

func (h Handler) HandlerSetlistPDFPrint(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistPDFPrint")

}

func (h Handler) HandlerSetlistPDFSave(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistPDFSave")
}
