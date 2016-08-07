package jsonfetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Fetcher retrives a body of text and marshalls it into the destination
type Fetcher interface {
	Fetch(string, interface{}) error
}

// Jsonfetcher is the struct wrapping the http client and response
type Jsonfetcher struct {
	client   *http.Client
	response *http.Response
}

// Fetch will create the http client, fetch the content and marshall it into the destination
func (jsonfetcher *Jsonfetcher) Fetch(url string, destination interface{}) error {
	jsonfetcher.createClient()

	if err := jsonfetcher.fetchResponse(url); err != nil {
		return err
	}

	defer jsonfetcher.response.Body.Close()

	if err := jsonfetcher.marshallResponse(destination); err != nil {
		return err
	}

	return nil
}

func (jsonfetcher *Jsonfetcher) createClient() {
	if jsonfetcher.client == nil {
		jsonfetcher.client = &http.Client{
			Timeout: time.Second * 10,
		}
	}
}

func (jsonfetcher *Jsonfetcher) fetchResponse(url string) error {
	var err error

	if jsonfetcher.response, err = jsonfetcher.client.Get(url); err != nil {
		return err
	}

	if jsonfetcher.response.StatusCode != 200 {
		err = fmt.Errorf("the request to %s returned with a non-200, %d", url, jsonfetcher.response.StatusCode)
	}

	return err
}

func (jsonfetcher *Jsonfetcher) marshallResponse(destination interface{}) error {
	var contents []byte
	var err error

	contents, err = ioutil.ReadAll(jsonfetcher.response.Body)

	err = json.Unmarshal(contents, destination)

	return err
}
