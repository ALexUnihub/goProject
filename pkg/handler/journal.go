package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	journal "github.com/ALexUnihub/goProject/pkg/jornal_repo"
)

type JournalHandler struct {
	JournalRepo journal.JournalRepo
}

func (h *JournalHandler) AddRecord(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, `unknown payload`, http.StatusBadRequest)
		return
	}

	newRecord, err := h.JournalRepo.AddItemInRepo(r)
	if err != nil {
		http.Error(w, "pckg handler, AddRecord, AddItemInRepo", http.StatusInternalServerError)
		log.Printf("pckg handler, AddRecord, err AddItemInRepo: %v", err.Error())
		return
	}

	byteValue, err := json.Marshal(newRecord)
	if err != nil {
		http.Error(w, "pckg handler, AddRecord, err in json marshal", http.StatusInternalServerError)
		log.Printf("pckg handler, AddRecord, err in json marshal: %v", err.Error())
		return
	}

	_, err = w.Write(byteValue)
	if err != nil {
		log.Printf("pckg handler, AddRecord, err Write byteValue in http ResponseWriter: %v", err.Error())
		return
	}
}

func (h *JournalHandler) GetGraph(w http.ResponseWriter, r *http.Request) {
	graph, err := h.JournalRepo.CreateGrapgh()
	if err != nil {
		http.Error(w, "pckg handler, GetGraph err", http.StatusInternalServerError)
		log.Printf("pckg handler, GetGraph err: %v", err.Error())
		return
	}

	byteValue, err := json.Marshal(graph)
	if err != nil {
		http.Error(w, "pckg handler, GetGraph err", http.StatusInternalServerError)
		log.Printf("pckg handler, GetGraph err: %v", err.Error())
		return
	}

	_, err = w.Write(byteValue)
	if err != nil {
		http.Error(w, "pckg handler, GetGraph err", http.StatusInternalServerError)
		log.Printf("pckg handler, GetGraph err: %v", err.Error())
		return
	}

	fmt.Println("  ", "Graph")
	for key, item := range graph.Nodes {
		fmt.Println(key, "\t", item)
	}

	fmt.Println("  ", "Table")
	for key, item := range graph.ActivityIndex {
		fmt.Println(key, "\t", item)
	}
}
