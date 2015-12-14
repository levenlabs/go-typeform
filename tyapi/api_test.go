package tyapi

import (
	"bytes"
	"encoding/json"
	"github.com/levenlabs/go-typeform/tyform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	. "testing"
)

type testClient struct {
	Body       io.ReadCloser
	StatusCode int
}

func (t *testClient) Do(r *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: t.StatusCode,
		Body:       t.Body,
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

	f := &tyform.Form{}
	err := json.Unmarshal(j, f)
	require.Nil(t, err)

	j = []byte(`{
		"id": "random",
		"urls": [{
			"id": "test",
			"form_id": "test1",
			"version": "v0.4"
		}]
	}`)
	client = &testClient{
		Body:       ioutil.NopCloser(bytes.NewBuffer(j)),
		StatusCode: http.StatusCreated,
	}

	res, err := Create(f)
	require.Nil(t, err)

	assert.Equal(t, "random", res.ID)
	assert.Equal(t, "test", res.URLs[0].ID)
	assert.Equal(t, "test1", res.URLs[0].FormID)
	assert.Equal(t, "v0.4", res.URLs[0].Version)
}
