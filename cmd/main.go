package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ALexUnihub/goProject/pkg/handler"
	journal "github.com/ALexUnihub/goProject/pkg/jornal_repo"
	"github.com/gorilla/mux"
)

func main() {
	journalRepo := journal.NewMemoryRepo()
	journalHandler := handler.JournalHandler{
		JournalRepo: journalRepo,
	}

	journal.AddStaticJournal(journalRepo)

	router := mux.NewRouter()

	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/api/record", journalHandler.AddRecord).Methods("POST")
	router.HandleFunc("/api/graph", journalHandler.GetGraph).Methods("GET")

	addr := ":8080"
	log.Println("server starting on addr :8080")
	http.ListenAndServe(addr, router)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	dataHTML, err := os.Open("../static/index.html")
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}

	byteValue, err := ioutil.ReadAll(dataHTML)
	if err != nil {
		http.Error(w, `Parsing html error`, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(byteValue)
	if err != nil {
		log.Println("err in pckg handlers, Index", err.Error())
		return
	}
}
