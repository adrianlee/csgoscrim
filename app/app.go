package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", RootHandler)

	// Static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello App")
}
