package jsonfetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Fetcher retrives a body of text and marshals it into the destination
type Fetcher interface {
	Get(string, interface{}) error
	LastResponseHeaders() map[string]string
}

// Jsonfetcher is the struct wrapping the http client and response
type Jsonfetcher struct {
	client   *http.Client
	response *http.Response
}

// LastResponseHeaders will return the headers from the last request
func (jsonfetcher *Jsonfetcher) LastResponseHeaders() map[string]string {
	var results = make(map[string]string)

	if jsonfetcher.response != nil {
		for header := range jsonfetcher.response.Header {
			results[header] = jsonfetcher.response.Header.Get(header)
		}
		return results
	}

	return nil
}

// Get will create the http client, fetch the content and marshal it into the destination
func (jsonfetcher *Jsonfetcher) Get(url string, headers map[string]string, destination interface{}) error {
	jsonfetcher.createClient()

	if err := jsonfetcher.get(url, headers); err != nil {
		return err
	}

	defer jsonfetcher.response.Body.Close()

	err := jsonfetcher.unmarshalResponse(destination)

	return err
}

func (jsonfetcher *Jsonfetcher) createClient() {
	if jsonfetcher.client == nil {
		jsonfetcher.client = &http.Client{
			Timeout: time.Second * 10,
		}
	}
}

func (jsonfetcher *Jsonfetcher) get(url string, headers map[string]string) error {
	var err error
	var request *http.Request

	if request, err = http.NewRequest("GET", url, nil); err != nil {
		return err
	}

	for header, value := range headers {
		request.Header.Set(header, value)
	}

	if jsonfetcher.response, err = jsonfetcher.client.Do(request); err != nil {
		return err
	}

	if jsonfetcher.response.StatusCode != 200 {
		err = fmt.Errorf("the request to %s returned with a non-200, %d", url, jsonfetcher.response.StatusCode)
	}

	return err
}

func (jsonfetcher *Jsonfetcher) unmarshalResponse(destination interface{}) error {
	var contents []byte
	var err error

	if contents, err = ioutil.ReadAll(jsonfetcher.response.Body); err != nil {
		return err
	}

	if err = json.Unmarshal(contents, destination); err != nil {
		return err
	}

	return nil
}
