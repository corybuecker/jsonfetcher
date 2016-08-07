package jsonfetcher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var server *httptest.Server

var expectedResponse = "{\"response\":{\"games\":[{\"appid\":10,\"playtime_forever\":32}]}}"
var unexpectedResponse = "{\"test\":{\"something\":true}}"

var fetcher = Jsonfetcher{}

type matchingDestination struct {
	Response struct {
		Games []struct {
			ID              int    `json:"appid"`
			Name            string `json:"name"`
			PlaytimeForever int    `json:"playtime_forever"`
			PlaytimeRecent  int    `json:"playtime_2weeks"`
		} `json:"games"`
	} `json:"response"`
}

type nonMatchingDestination struct {
	Response struct {
	} `json:"bad"`
}

func configureResponse(code int, response string) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		fmt.Fprint(w, response)
	}))
}

func TestMatchingDestinationWithExpectedResponse(t *testing.T) {
	configureResponse(200, expectedResponse)
	var data = matchingDestination{}
	fetcher.Fetch(server.URL, &data)

	assert.Equal(t, 10, data.Response.Games[0].ID, "should have returned the correct data")
}
