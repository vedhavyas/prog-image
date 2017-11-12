package progimg

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// StartImageServer will start the image server
func StartImageServer(addr string) {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handle404)
	r.HandleFunc("/images", handleUpload).Methods("POST")
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("failed to start server: %v\n", err)
	}
}
