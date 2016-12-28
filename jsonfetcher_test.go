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
var malformedResponse = "\"test\":{\"something\":true}}"

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

func TestHTTPError(t *testing.T) {
	configureResponse(200, expectedResponse)
	var data = matchingDestination{}
	server.Close()

	err := fetcher.Get(server.URL, nil, &data)

	assert.Error(t, err, "should be an error")
}

func TestMatchingDestinationWithExpectedResponse(t *testing.T) {
	configureResponse(200, expectedResponse)
	var data = matchingDestination{}
	fetcher.Get(server.URL, nil, &data)

	assert.Equal(t, 10, data.Response.Games[0].ID, "should have returned the correct data")
}

func TestNonMatchingDestinationWithExpectedResponse(t *testing.T) {
	configureResponse(200, expectedResponse)
	var data = nonMatchingDestination{}
	fetcher.Get(server.URL, nil, &data)

	assert.Equal(t, nonMatchingDestination{}, data, "should be empty")
}

func TestMatchingDestinationWithUnexpectedResponse(t *testing.T) {
	configureResponse(200, unexpectedResponse)
	var data = matchingDestination{}
	fetcher.Get(server.URL, nil, &data)

	assert.Equal(t, matchingDestination{}, data, "should be empty")
}

func TestNonMatchingDestinationWithUnexpectedResponse(t *testing.T) {
	configureResponse(200, unexpectedResponse)
	var data = nonMatchingDestination{}
	fetcher.Get(server.URL, nil, &data)

	assert.Equal(t, nonMatchingDestination{}, data, "should be empty")
}

func TestNonTwoHundredResponseCode(t *testing.T) {
	configureResponse(500, expectedResponse)
	var data = matchingDestination{}
	err := fetcher.Get(server.URL, nil, &data)
	assert.Error(t, err, "should be an error")
}

func TestMalformedResponse(t *testing.T) {
	configureResponse(200, malformedResponse)
	var data = matchingDestination{}
	err := fetcher.Get(server.URL, nil, &data)
	assert.Error(t, err, "should be an error")
}

func TestHeaders(t *testing.T) {
	configureResponse(200, expectedResponse)
	var data = matchingDestination{}
	fetcher.Get(server.URL, map[string]string{"testheader": "true"}, &data)
	assert.Equal(t, 10, data.Response.Games[0].ID, "should have returned the correct data")
}
