package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "testing"
	"encoding/json"
	"github.com/levenlabs/go-typeform/fields"
	"net/http"
	"io"
	"io/ioutil"
	"bytes"
)

type testClient struct {
	Body io.ReadCloser
	StatusCode int
}

func (t *testClient) Do(r *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: t.StatusCode,
		Body: t.Body,
	}
	return resp, nil
}

func init() {
	APIToken = "test"
}

func TestCreate(t *T) {
	j := []byte(`{
		"title": "Form",
		"fields": [
			{
				"type": "statement",
				"question": "Hey"
			}
		]
	}`)

	f := &fields.Form{}
	err := json.Unmarshal(j, f)
	require.Nil(t, err)

	j = []byte(`{
		"id": "random",
		"urls": {
			"id": "test",
			"form_id": "test1",
			"version": "v0.4"
		}
	}`)
	client = &testClient{
		Body: ioutil.NopCloser(bytes.NewBuffer(j)),
		StatusCode: http.StatusCreated,
	}

	res, err := Create(f)
	require.Nil(t, err)

	assert.Equal(t, "random", res.ID)
	assert.Equal(t, "test", res.URLs.ID)
	assert.Equal(t, "test1", res.URLs.FormID)
	assert.Equal(t, "v0.4", res.URLs.Version)
}