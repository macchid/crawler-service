package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/macchid/crawler-service/hfuncs"
	"github.com/macchid/crawler-service/middleware"
)

/*
	This will be a little example service that will have 5 simple endpoints

	- POST /api/v1/crawl: Creates and returns a crawl's result in JSON format.
	- GET /api/v1/crawl: Return all crawls that are currently running made in JSON format.
	- GET /api/v1/crawl/{id}: Check's crawl with {id} status and get result if exists.
	- DELETE /api/v1/crawl/{id}: Forcibly stops crawl with ID {id}.

	- PATCH /api/v1/crawl/{id}: Re-runs and update results for crawl with ID {id}.

*/

func apiHandlers(root *mux.Router) {
	api := root.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/crawl", hfuncs.NewCrawlHandler).Methods(http.MethodPost)
	api.HandleFunc("/crawl", hfuncs.ListCrawlsHandler).Methods(http.MethodGet)
	api.HandleFunc("/crawl/{id}", hfuncs.RepeatCrawlHandler).Methods(http.MethodGet)
	api.HandleFunc("/crawl/{id}", hfuncs.RepeatCrawlHandler).Methods(http.MethodPatch)
	api.HandleFunc("/crawl/{id}", hfuncs.StopCrawlHandler).Methods(http.MethodDelete)
}

func main() {
	// TODO: Read configuration file, or use flags. (go/pkg/flags)

	root := mux.NewRouter()

	f, err := os.OpenFile("logs/crawler.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Couldn't open the log file")
		return
	}
	defer f.Close()

	logger := log.New(f, "[HANDLER]", log.LstdFlags|log.Lmicroseconds|log.Lmsgprefix|log.Llongfile)

	root.Use(middleware.LoggingMiddleware(logger))

	root.HandleFunc("/", hfuncs.HomeHandler)
	apiHandlers(root)

	log.Fatal(http.ListenAndServe(":8080", root))
}
