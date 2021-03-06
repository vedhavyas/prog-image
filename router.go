package progimg

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// getRouter returns a mux router with defined urls
func getRouter() http.Handler {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handle404)
	r.HandleFunc("/images/{id}", handleDownload).Methods("GET")
	r.HandleFunc("/images/", handleUpload).Methods("POST")
	r.HandleFunc("/images", handleUpload).Methods("POST")
	return r
}

// StartImageServer will start the image server
func StartImageServer(addr string) {
	err := http.ListenAndServe(addr, recoverHandler(logHandler(getRouter())))
	if err != nil {
		log.Fatalf("failed to start server: %v\n", err)
	}
}
