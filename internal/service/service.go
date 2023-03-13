package service

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func NewService() {

	r := mux.NewRouter()

	r.HandleFunc("/", handleConnection)

	//http.Handle("/", r)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln(err)
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "Category: %v\n", vars["category"])
	w.Write([]byte("OK"))
}
