package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/levenlabs/go-typeform/fields"
	"net/http"
)

var version = "v0.4"
var APIToken string
var errRmptyToken = errors.New("Empty APIToken")
var client httpClient = http.DefaultClient

type URLs struct {
	ID      string `json:"id"`
	FormID  string `json:"form_id"`
	Version string `json:"version"`
}

type CreateResult struct {
	ID   string `json:"id"`
	URLs []URLs `json:"urls"`
}

func Create(f *fields.Form) (*CreateResult, error) {
	if APIToken == "" {
		return nil, errRmptyToken
	}

	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("https://api.typeform.io/%s/forms", version)
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-TOKEN", APIToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected response from /forms: %s", resp.Status)
	}

	res := &CreateResult{}
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(res); err != nil {
		return nil, err
	}
	return res, nil
}

// httpClient is an interface that describes http.Client so we can override what
// client we use in testing
type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}
