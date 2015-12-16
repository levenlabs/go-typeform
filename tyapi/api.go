package tyapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/levenlabs/go-typeform/tyform"
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

// Creates a survey on typeform. Returns an `Error` if we get one.
func Create(f *tyform.Form) (*CreateResult, error) {
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
		errRes := &Error{}
		dec := json.NewDecoder(resp.Body)
		if err = dec.Decode(errRes); err != nil {
			return nil, fmt.Errorf("unexpected response from /forms: %s", resp.Status)
		}
		return nil, *errRes
	}

	res := &CreateResult{}
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(res); err != nil {
		return nil, err
	}
	return res, nil
}

type Error struct {
	ErrorType   string `json:"error"`
	Field       string `json:"field"`
	Description string `json:"description"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s on field %s: %s", e.ErrorType, e.Field, e.Description)
}

// httpClient is an interface that describes http.Client so we can override what
// client we use in testing
type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}
