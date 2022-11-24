package deribit

import (
	"bytes"
	"encoding/json"
	"github.com/DawnKosmos/metapine/backend/exchange"
	"io"
	"log"
	"net/http"
)

type Deribit struct {
	client *http.Client
}

const URL = "https://deribit.com/api/v2/"

func New() exchange.CandleProvider {
	return &Deribit{client: http.DefaultClient}
}

func (d *Deribit) Name() string {
	return "Deribit"
}

func (d *Deribit) request(method string, path string, body []byte) *http.Request {
	req, _ := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func (d *Deribit) get(path string, body []byte) (*http.Response, error) {
	preparedRequest := d.request("GET", path, body)
	resp, err := d.client.Do(preparedRequest)
	return resp, err
}

func processResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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
