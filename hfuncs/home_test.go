package hfuncs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomeHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HomeHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code: got [%v] want [%v]", rr.Code, http.StatusOK)

	expected := "Welcome to the Crawler Service"
	assert.Equal(t, rr.Body.String(), expected, "Handler returned unexpected body: got [%v] want [%v]", rr.Body.String(), expected)
}
