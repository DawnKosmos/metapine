package ftx

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type FTX struct {
	Client *http.Client
}

func (f *FTX) Name() string {
	return "ftx"
}

const URL = "https://ftx.com/api/"

func New() *FTX {
	return &FTX{&http.Client{}}
}

func (f *FTX) request(method string, path string, body []byte) *http.Request {
	req, _ := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func (f *FTX) get(path string, body []byte) (*http.Response, error) {
	preparedRequest := f.request("GET", path, body)
	resp, err := f.Client.Do(preparedRequest)
	return resp, err
}

func processResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error processing response: %v", err)
		return err
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		log.Printf("Error processing response: %v", err)
		return err
	}
	return nil
}
