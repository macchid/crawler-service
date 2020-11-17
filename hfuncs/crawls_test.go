package hfuncs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCrawlHandler(t *testing.T) {
	// Create the body
	reader := strings.NewReader("{\"url\":\"https://golang.org/\", \"depth\":5}")

	// Create the request
	req, err := http.NewRequest("POST", "/api/v1/crawl", reader)
	if err != nil {
		t.Fatal(err)
	}

	// Set the content-type header
	req.Header.Add("Content-Type", "application/json")

	// Create a request recorder
	rr := httptest.NewRecorder()

	// Define the handler and Serve
	handler := http.HandlerFunc(NewCrawlHandler)
	handler.ServeHTTP(rr, req)

	// Assert status code 200 (OK)
	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code: got [%v] want [%v]", rr.Code, http.StatusOK)

	// Assert body contains "The Go Programming Language"
	mustContain := "The Go Programming Language"
	assert.Contains(t, rr.Body.String(), mustContain, "Hanlder returned unexpected body: got [%v] wanted to contain [%v]", rr.Body.String(), mustContain)
}

func Test