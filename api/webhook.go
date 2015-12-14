package api

import (
	"encoding/json"
	"github.com/levenlabs/go-llog"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

// Results represents a single set of answers received from typeform via the
// webhook in response to someone taking a form
type Results struct {
	ID      string          `json:"id"              bson:"i"`
	Token   string          `json:"token"           bson:"t"`
	Answers []ResultsAnswer `json:"answers"         bson:"a"`
}

// ResultsAnswerMetadata is shared between the different forms of the answer
type ResultsAnswerMetadata struct {
	FieldID int64    `json:"field_id" bson:"i"`
	Type    string   `json:"type" bson:"t"`
	Tags    []string `json:"tags,omitempty"         bson:"g,omitempty"`
}

// ResultsAnswer is a single answer in the Results
// The Value can be many different types. See http://docs.typeform.io/docs/results-introduction
type ResultsAnswer struct {
	ResultsAnswerMetadata `bson:",inline"`
	Value                 interface{} `json:"value" bson:"v"`
}

// jsonForm is used to Unmarshal into since it has Value of json.RawMessage
type jsonAnswer struct {
	ResultsAnswerMetadata
	Value json.RawMessage `json:"value"             bson:"v"`
}

// bsonForm is used to Unmarshal into since it has Value of bson.Raw
type bsonAnswer struct {
	ResultsAnswerMetadata `bson:",inline"`
	Value                 bson.Raw `json:"value"    bson:"v"`
}

// NumberValue represents a number answer
type NumberValue struct {
	Amount int64 `json:"amount"                     bson:"a"`
}

// ChoiceValue represents a number answer
type ChoiceValue struct {
	Label      string `json:"label"                 bson:"l,omitempty"`
	Other      string `json:"other,omitempty"       bson:"o,omitempty"`
	EmptyOther bool   `json:"-"                     bson:"eo,omitempty"`
}

type jsonChoiceValue struct {
	Label string          `json:"label"             bson:"l,omitempty"`
	Other json.RawMessage `json:"other,omitempty"`
}

type ChoicesValue struct {
	Labels     []string `json:"labels"              bson:"l,omitempty"`
	Other      string   `json:"other,omitempty"     bson:"o,omitempty"`
	EmptyOther bool     `json:"-"                   bson:"eo,omitempty"`
}

type jsonChoicesValue struct {
	Labels []string        `json:"labels"           bson:"l,omitempty"`
	Other  json.RawMessage `json:"other,omitempty"`
}

type TextValue string

type BooleanValue bool

// ListenAndServe starts an http server at the given addr and requires a handler
// that will be called for each webhook request and passed a Results pointer. If
// an error is returned from the handler, a 500 response is sent and TypeForm
// will retry the request.
//
// Note: the handler might be called multiple times for the same results set.
// You should store the Token for each call and verify you haven't already
// processed it.
func ListenAndServe(addr string, cb func(*Results) error) error {
	return http.ListenAndServe(addr, http.HandlerFunc(wrapCallback(cb)))
}

func wrapCallback(cb func(*Results) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		kv := llog.KV{
			"ip":  r.RemoteAddr,
			"url": r.URL.String(),
		}
		if r.Method != "POST" {
			kv["method"] = r.Method
			llog.Warn("invalid method received at webhook", kv)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		d := json.NewDecoder(r.Body)
		res := &Results{}
		if err := d.Decode(res); err != nil {
			kv["error"] = err
			llog.Warn("json error while decoding webhook body", kv)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := cb(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (a *ResultsAnswer) emptyValue() interface{} {
	switch a.Type {
	case "number":
		return &NumberValue{}
	case "choice":
		return &ChoiceValue{}
	case "choices":
		return &ChoicesValue{}
	case "text":
		return new(TextValue)
	case "boolean":
		return new(BooleanValue)
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (a *ResultsAnswer) UnmarshalJSON(b []byte) error {
	ja := &jsonAnswer{}
	if err := json.Unmarshal(b, ja); err != nil {
		return err
	}
	a.ResultsAnswerMetadata = ja.ResultsAnswerMetadata

	a.Value = a.emptyValue()
	if a.Value != nil {
		if err := json.Unmarshal(ja.Value, a.Value); err != nil {
			return err
		}
	}
	return nil
}

// SetBSON implements the bson.Setter interface
func (a *ResultsAnswer) SetBSON(raw bson.Raw) error {
	ba := &bsonAnswer{}
	if err := raw.Unmarshal(ba); err != nil {
		return err
	}
	a.ResultsAnswerMetadata = ba.ResultsAnswerMetadata

	a.Value = a.emptyValue()
	if a.Value != nil {
		if err := ba.Value.Unmarshal(a.Value); err != nil {
			return err
		}
	}
	return nil
}

func otherUnmarshal(m json.RawMessage, dst *string) (bool, error) {
	if len(m) == 0 || string(m) == "null" {
		return true, nil
	}
	if serr := json.Unmarshal(m, dst); serr != nil {
		n := new(int64)
		if ierr := json.Unmarshal(m, n); ierr != nil {
			return false, serr
		}
		*dst = strconv.FormatInt(*n, 10)
	}
	return false, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *ChoiceValue) UnmarshalJSON(b []byte) error {
	jv := &jsonChoiceValue{}
	var err error
	if err := json.Unmarshal(b, jv); err != nil {
		return err
	}
	v.Label = jv.Label
	v.EmptyOther, err = otherUnmarshal(jv.Other, &v.Other)
	return err
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *ChoicesValue) UnmarshalJSON(b []byte) error {
	jv := &jsonChoicesValue{}
	var err error
	if err := json.Unmarshal(b, jv); err != nil {
		return err
	}
	v.Labels = jv.Labels
	v.EmptyOther, err = otherUnmarshal(jv.Other, &v.Other)
	return err
}
