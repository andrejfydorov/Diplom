package service

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func NewService() {
	r := mux.NewRouter()

	r.HandleFunc("/", handleConnection)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println(err)
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
