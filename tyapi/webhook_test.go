package tyapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	. "testing"
)

func TestJSONNumber(t *T) {
	a := &ResultsAnswer{
		Value: &NumberValue{
			Amount: 5,
		},
	}
	a.Type = "number"
	fs := `{"field_id":0,"type":"number","value":{"amount":5}}`
	j, err := json.Marshal(a)
	require.Nil(t, err)
	assert.Equal(t, fs, string(j))

	na := &ResultsAnswer{}
	err = json.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)
}

func TestJSONNumberExp(t *T) {
	a := &ResultsAnswer{
		ResultsAnswerMetadata: ResultsAnswerMetadata{
			FieldID: 1774766,
		},
		Value: &NumberValue{
			Amount: 1111111,
		},
	}
	a.Type = "number"
	fs := []byte(`{"field_id":1.774766e+06,"type":"number","value":{"amount":1.111111e+06}}`)
	na := &ResultsAnswer{}
	err := json.Unmarshal(fs, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)
}

func TestBSONNumber(t *T) {
	n := &NumberValue{
		Amount: 5,
	}
	a := &ResultsAnswer{
		Value: n,
	}
	a.Type = "number"
	aexp := struct {
		FieldID int64       `bson:"i"`
		Type    string      `bson:"t"`
		Value   interface{} `bson:"v"`
	}{
		Type:  "number",
		Value: n,
	}
	j, err := bson.Marshal(a)
	require.Nil(t, err)
	jexp, err := bson.Marshal(aexp)
	require.Nil(t, err)
	assert.Equal(t, string(jexp), string(j))

	na := &ResultsAnswer{}
	err = bson.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)
}

func TestJSONBoolean(t *T) {
	v := BooleanValue(true)
	a := &ResultsAnswer{
		Value: &v,
	}
	a.Type = "boolean"
	fs := `{"field_id":0,"type":"boolean","value":true}`
	j, err := json.Marshal(a)
	require.Nil(t, err)
	assert.Equal(t, fs, string(j))

	na := &ResultsAnswer{}
	err = json.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)

	fs = `{"field_id":0,"type":"boolean","value":false}`
	na = &ResultsAnswer{}
	err = json.Unmarshal([]byte(fs), na)
	require.Nil(t, err)
	v = BooleanValue(false)
	assert.Equal(t, &v, na.Value)
}

func TestBSONBoolean(t *T) {
	v := BooleanValue(true)
	a := &ResultsAnswer{
		Value: &v,
	}
	a.Type = "boolean"
	aexp := struct {
		FieldID int64       `bson:"i"`
		Type    string      `bson:"t"`
		Value   interface{} `bson:"v"`
	}{
		Type:  "boolean",
		Value: &v,
	}
	j, err := bson.Marshal(a)
	require.Nil(t, err)
	jexp, err := bson.Marshal(aexp)
	require.Nil(t, err)
	assert.Equal(t, string(jexp), string(j))

	na := &ResultsAnswer{}
	err = bson.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)
}

func TestJSONText(t *T) {
	v := TextValue("hey")
	a := &ResultsAnswer{
		Value: &v,
	}
	a.Type = "text"
	fs := `{"field_id":0,"type":"text","value":"hey"}`
	j, err := json.Marshal(a)
	require.Nil(t, err)
	assert.Equal(t, fs, string(j))

	na := &ResultsAnswer{}
	err = json.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)

	fs = `{"field_id":0,"type":"text","value":""}`
	na = &ResultsAnswer{}
	err = json.Unmarshal([]byte(fs), na)
	require.Nil(t, err)
	v = TextValue("")
	assert.Equal(t, &v, na.Value)
}

func TestBSONText(t *T) {
	v := TextValue("hey")
	a := &ResultsAnswer{
		Value: &v,
	}
	a.Type = "text"
	aexp := struct {
		FieldID int64       `bson:"i"`
		Type    string      `bson:"t"`
		Value   interface{} `bson:"v"`
	}{
		Type:  "text",
		Value: &v,
	}
	j, err := bson.Marshal(a)
	require.Nil(t, err)
	jexp, err := bson.Marshal(aexp)
	require.Nil(t, err)
	assert.Equal(t, string(jexp), string(j))

	na := &ResultsAnswer{}
	err = bson.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)
}

func TestJSONChoice(t *T) {
	a := &ResultsAnswer{
		Value: &ChoiceValue{
			Label:      "val",
			EmptyOther: true,
		},
	}
	a.Type = "choice"
	fs := `{"field_id":0,"type":"choice","value":{"label":"val"}}`
	j, err := json.Marshal(a)
	require.Nil(t, err)
	assert.Equal(t, fs, string(j))

	na := &ResultsAnswer{}
	err = json.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)

	a = &ResultsAnswer{
		Value: &ChoiceValue{
			Other: "o",
		},
	}
	a.Type = "choice"
	fs = `{"field_id":0,"type":"choice","value":{"label":"","other":"o"}}`
	j, err = json.Marshal(a)
	require.Nil(t, err)
	assert.Equal(t, fs, string(j))

	na = &ResultsAnswer{}
	err = json.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)

	fs = `{"field_id":0,"type":"choice","value":{"other":1}}`
	na = &ResultsAnswer{}
	err = json.Unmarshal([]byte(fs), na)
	require.Nil(t, err)
	v, ok := na.Value.(*ChoiceValue)
	require.True(t, ok)
	assert.Equal(t, "1", v.Other)

	fs = `{"field_id":0,"type":"choice","value":{"other":null}}`
	na = &ResultsAnswer{}
	err = json.Unmarshal([]byte(fs), na)
	require.Nil(t, err)
	v, ok = na.Value.(*ChoiceValue)
	require.True(t, ok)
	assert.Equal(t, "", v.Other)
	assert.True(t, v.EmptyOther)

	fs = `{"field_id":0,"type":"choice","value":{}}`
	na = &ResultsAnswer{}
	err = json.Unmarshal([]byte(fs), na)
	require.Nil(t, err)
	v, ok = na.Value.(*ChoiceValue)
	require.True(t, ok)
	assert.Equal(t, "", v.Other)
	assert.True(t, v.EmptyOther)
}

func TestBSONChoice(t *T) {
	n := &ChoiceValue{
		Label: "val",
	}
	a := &ResultsAnswer{
		Value: n,
	}
	a.Type = "choice"
	aexp := struct {
		FieldID int64       `bson:"i"`
		Type    string      `bson:"t"`
		Value   interface{} `bson:"v"`
	}{
		Type:  "choice",
		Value: n,
	}
	j, err := bson.Marshal(a)
	require.Nil(t, err)
	jexp, err := bson.Marshal(aexp)
	require.Nil(t, err)
	assert.Equal(t, string(jexp), string(j))

	na := &ResultsAnswer{}
	err = bson.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)

	n = &ChoiceValue{
		Other: "o",
	}
	a = &ResultsAnswer{
		Value: n,
	}
	a.Type = "choice"
	aexp = struct {
		FieldID int64       `bson:"i"`
		Type    string      `bson:"t"`
		Value   interface{} `bson:"v"`
	}{
		Type:  "choice",
		Value: n,
	}
	j, err = bson.Marshal(a)
	require.Nil(t, err)
	jexp, err = bson.Marshal(aexp)
	require.Nil(t, err)
	assert.Equal(t, string(jexp), string(j))

	na = &ResultsAnswer{}
	err = bson.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)

	n = &ChoiceValue{
		EmptyOther: true,
	}
	a = &ResultsAnswer{
		Value: n,
	}
	a.Type = "choice"
	aexp = struct {
		FieldID int64       `bson:"i"`
		Type    string      `bson:"t"`
		Value   interface{} `bson:"v"`
	}{
		Type:  "choice",
		Value: n,
	}
	j, err = bson.Marshal(a)
	require.Nil(t, err)
	jexp, err = bson.Marshal(aexp)
	require.Nil(t, err)
	assert.Equal(t, string(jexp), string(j))

	na = &ResultsAnswer{}
	err = bson.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)
}

func TestJSONChoices(t *T) {
	a := &ResultsAnswer{
		Value: &ChoicesValue{
			Labels:     []string{"val"},
			EmptyOther: true,
		},
	}
	a.Type = "choices"
	fs := `{"field_id":0,"type":"choices","value":{"labels":["val"]}}`
	j, err := json.Marshal(a)
	require.Nil(t, err)
	assert.Equal(t, fs, string(j))

	na := &ResultsAnswer{}
	err = json.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)

	a = &ResultsAnswer{
		Value: &ChoicesValue{
			Other: "o",
		},
	}
	a.Type = "choices"
	fs = `{"field_id":0,"type":"choices","value":{"labels":null,"other":"o"}}`
	j, err = json.Marshal(a)
	require.Nil(t, err)
	assert.Equal(t, fs, string(j))

	na = &ResultsAnswer{}
	err = json.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)

	fs = `{"field_id":0,"type":"choices","value":{"other":1}}`
	na = &ResultsAnswer{}
	err = json.Unmarshal([]byte(fs), na)
	require.Nil(t, err)
	v, ok := na.Value.(*ChoicesValue)
	require.True(t, ok)
	assert.Equal(t, "1", v.Other)

	fs = `{"field_id":0,"type":"choices","value":{"other":null}}`
	na = &ResultsAnswer{}
	err = json.Unmarshal([]byte(fs), na)
	require.Nil(t, err)
	v, ok = na.Value.(*ChoicesValue)
	require.True(t, ok)
	assert.Equal(t, "", v.Other)
	assert.True(t, v.EmptyOther)

	fs = `{"field_id":0,"type":"choices","value":{}}`
	na = &ResultsAnswer{}
	err = json.Unmarshal([]byte(fs), na)
	require.Nil(t, err)
	v, ok = na.Value.(*ChoicesValue)
	require.True(t, ok)
	assert.Equal(t, "", v.Other)
	assert.True(t, v.EmptyOther)
}

func TestBSONChoices(t *T) {
	n := &ChoicesValue{
		Labels: []string{"val"},
	}
	a := &ResultsAnswer{
		Value: n,
	}
	a.Type = "choices"
	aexp := struct {
		FieldID int64       `bson:"i"`
		Type    string      `bson:"t"`
		Value   interface{} `bson:"v"`
	}{
		Type:  "choices",
		Value: n,
	}
	j, err := bson.Marshal(a)
	require.Nil(t, err)
	jexp, err := bson.Marshal(aexp)
	require.Nil(t, err)
	assert.Equal(t, string(jexp), string(j))

	na := &ResultsAnswer{}
	err = bson.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)

	n = &ChoicesValue{
		Other: "o",
	}
	a = &ResultsAnswer{
		Value: n,
	}
	a.Type = "choices"
	aexp = struct {
		FieldID int64       `bson:"i"`
		Type    string      `bson:"t"`
		Value   interface{} `bson:"v"`
	}{
		Type:  "choices",
		Value: n,
	}
	j, err = bson.Marshal(a)
	require.Nil(t, err)
	jexp, err = bson.Marshal(aexp)
	require.Nil(t, err)
	assert.Equal(t, string(jexp), string(j))

	na = &ResultsAnswer{}
	err = bson.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)

	n = &ChoicesValue{
		EmptyOther: true,
	}
	a = &ResultsAnswer{
		Value: n,
	}
	a.Type = "choices"
	aexp = struct {
		FieldID int64       `bson:"i"`
		Type    string      `bson:"t"`
		Value   interface{} `bson:"v"`
	}{
		Type:  "choices",
		Value: n,
	}
	j, err = bson.Marshal(a)
	require.Nil(t, err)
	jexp, err = bson.Marshal(aexp)
	require.Nil(t, err)
	assert.Equal(t, string(jexp), string(j))

	na = &ResultsAnswer{}
	err = bson.Unmarshal(j, na)
	require.Nil(t, err)
	assert.EqualValues(t, a, na)
}

func TestWrapCallback(t *T) {
	b := []byte(`{
		"uid": "test",
		"token": "t1",
		"answers": [
			{"field_id":123,"type":"boolean","value":true},
			{"field_id":122,"type":"boolean","value":true}
		]
	}`)
	u, _ := url.Parse("http://test")
	r := httptest.NewRecorder()
	req := &http.Request{
		Method: "POST",
		Body:   ioutil.NopCloser(bytes.NewBuffer(b)),
		URL:    u,
	}
	wrapCallback(func(r *Results, _ *http.Request) error {
		assert.Equal(t, "test", r.UID)
		assert.Equal(t, "t1", r.Token)
		require.Len(t, r.Answers, 2)
		assert.EqualValues(t, 122, r.Answers[0].FieldID)
		assert.Equal(t, "boolean", r.Answers[0].Type)
		assert.EqualValues(t, 123, r.Answers[1].FieldID)
		v := BooleanValue(true)
		assert.Equal(t, &v, r.Answers[0].Value)
		return nil
	})(r, req)
	assert.Equal(t, http.StatusOK, r.Code)

	b = []byte(`{,}`)
	r = httptest.NewRecorder()
	req = &http.Request{
		Method: "POST",
		Body:   ioutil.NopCloser(bytes.NewBuffer(b)),
		URL:    u,
	}
	wrapCallback(func(r *Results, _ *http.Request) error {
		// this should never run
		require.True(t, false)
		return nil
	})(r, req)
	assert.Equal(t, http.StatusBadRequest, r.Code)

	r = httptest.NewRecorder()
	req = &http.Request{
		Method: "GET",
		URL:    u,
	}
	wrapCallback(func(r *Results, _ *http.Request) error {
		// this should never run
		require.True(t, false)
		return nil
	})(r, req)
	assert.Equal(t, http.StatusMethodNotAllowed, r.Code)

	b = []byte(`{}`)
	r = httptest.NewRecorder()
	req = &http.Request{
		Method: "POST",
		Body:   ioutil.NopCloser(bytes.NewBuffer(b)),
		URL:    u,
	}
	wrapCallback(func(r *Results, _ *http.Request) error {
		return errors.New("an error occurred")
	})(r, req)
	assert.Equal(t, http.StatusInternalServerError, r.Code)
	assert.Equal(t, 0, r.Body.Len())
}

func TestAnswersSort(t *T) {
	s := ResultsAnswerSlice{
		ResultsAnswer{
			ResultsAnswerMetadata: ResultsAnswerMetadata{
				FieldID: 9,
			},
		},
		ResultsAnswer{
			ResultsAnswerMetadata: ResultsAnswerMetadata{
				FieldID: 1,
			},
		},
		ResultsAnswer{
			ResultsAnswerMetadata: ResultsAnswerMetadata{
				FieldID: 3,
			},
		},
	}
	sort.Sort(s)
	assert.EqualValues(t, 1, s[0].FieldID)
	assert.EqualValues(t, 3, s[1].FieldID)
	assert.EqualValues(t, 9, s[2].FieldID)
}

func TestAnswerString(t *T) {
	a := &ResultsAnswer{
		Value: &NumberValue{
			Amount: 5,
		},
	}
	assert.Equal(t, "5", a.String())

	b := BooleanValue(true)
	a = &ResultsAnswer{
		Value: &b,
	}
	assert.Equal(t, "true", a.String())
	b = BooleanValue(false)
	a = &ResultsAnswer{
		Value: &b,
	}
	assert.Equal(t, "false", a.String())

	tv := TextValue("hey")
	a = &ResultsAnswer{
		Value: &tv,
	}
	assert.Equal(t, "hey", a.String())

	a = &ResultsAnswer{
		Value: &ChoiceValue{
			Label:      "val",
			EmptyOther: true,
		},
	}
	assert.Equal(t, "val", a.String())
	a = &ResultsAnswer{
		Value: &ChoiceValue{
			Other:      "other",
			EmptyOther: false,
		},
	}
	assert.Equal(t, "other", a.String())

	a = &ResultsAnswer{
		Value: &ChoicesValue{
			Labels:     []string{"val", "val2"},
			EmptyOther: true,
		},
	}
	assert.Equal(t, "val,val2", a.String())
	a = &ResultsAnswer{
		Value: &ChoicesValue{
			Other:      "other",
			EmptyOther: false,
		},
	}
	assert.Equal(t, "other", a.String())
}
