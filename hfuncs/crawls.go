package hfuncs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/macchid/crawler-service/crawler"
)

type crawlParams struct {
	URL   string `json:"url"`
	Depth int    `json:"depth"`
}

type crawlError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type crawlResult struct {
	ID      string         `json:"id"`
	Params  crawlParams    `json:"params"`
	Results []crawler.Item `json:"results"`
}

// Retrieve a list of all Crawls perfored until now.
func ListCrawlsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Pending implementation for endpoint %v", r.RequestURI)
}

func NewCrawlHandler(w http.ResponseWriter, r *http.Request) {
	// Get the crawler from the request Body
	newCrawl := crawlResult{}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errjson := crawlError{
			http.StatusBadRequest,
			"Kindly enter data with url and depth to be crawled",
		}
		json.NewEncoder(w).Encode(errjson)
	}
	json.Unmarshal(reqBody, &(newCrawl.Params))

	dispatched := crawler.Dispatch(newCrawl.Params.URL, newCrawl.Params.Depth)
	newCrawl.ID = dispatched.ID()

	t := time.AfterFunc(10*time.Second, func() {
		log.Print("Calling the 10 seconds timer")
		err := dispatched.Close()
		json.NewEncoder(w).Encode(crawlError{Status: http.StatusGatewayTimeout, Message: fmt.Sprint(err)})
	})
	defer t.Stop()

	for it := range dispatched.Get() {
		newCrawl.Results = append(newCrawl.Results, it)
	}

	json.NewEncoder(w).Encode(newCrawl)
}

func GetCrawlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	msg := crawlError{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("You requested id %v. In the future you'll get results instead of this message", id),
	}

	json.NewEncoder(w).Encode(msg)
}

func RepeatCrawlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	newCrawl := crawlResult{ID: id}
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(crawlError{Status: http.StatusInternalServerError, Message: fmt.Sprint(err)})
	}

	json.Unmarshal(content, &(newCrawl.Params))

	json.NewEncoder(w).Encode(
		crawlError{
			Status:  http.StatusOK,
			Message: fmt.Sprintf("You requested to run again %v with this params: {url:%v, depth:%v}. In the future, this will start a new crawl.", newCrawl.ID, newCrawl.Params.URL, newCrawl.Params.Depth),
		},
	)
}

func StopCrawlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	json.NewEncoder(w).Encode(
		crawlError{
			Status:  http.StatusOK,
			Message: fmt.Sprintf("You requested to stop %v. In the future, this will end with the crawler stopped and with it's results deleted from storage", id),
		},
	)
}
