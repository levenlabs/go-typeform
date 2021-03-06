package tyapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/levenlabs/go-llog"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// ResultsAnswerSlice implements the sort interface and sorts by FieldID
// ascending
type ResultsAnswerSlice []ResultsAnswer

// Results represents a single set of answers received from typeform via the
// webhook in response to someone taking a form
type Results struct {
	UID     string             `json:"uid"          bson:"i"`
	Token   string             `json:"token"        bson:"t"`
	Answers ResultsAnswerSlice `json:"answers"      bson:"a"`
}

// ResultsAnswerMetadata is shared between the different forms of the answer
type ResultsAnswerMetadata struct {
	FieldID int64    `json:"field_id" bson:"i"`
	Type    string   `json:"type" bson:"t"`
	Tags    []string `json:"tags,omitempty"         bson:"g,omitempty"`
}

// jsonResultsAnswerMetadata is needed because json has to use json.Number
type jsonResultsAnswerMetadata struct {
	FieldID json.Number `json:"field_id" bson:"i"`
	Type    string      `json:"type" bson:"t"`
	Tags    []string    `json:"tags,omitempty"      bson:"g,omitempty"`
}

// ResultsAnswer is a single answer in the Results
// The Value can be many different types. See http://docs.typeform.io/docs/results-introduction
type ResultsAnswer struct {
	ResultsAnswerMetadata `bson:",inline"`
	Value                 interface{} `json:"value" bson:"v"`
}

// jsonForm is used to Unmarshal into since it has Value of json.RawMessage
type jsonAnswer struct {
	jsonResultsAnswerMetadata
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

// jsonNumberValue is needed because json has to use json.Number
type jsonNumberValue struct {
	Amount json.Number `json:"amount"`
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
func ListenAndServe(addr string, cb func(*Results, *http.Request) error) error {
	return http.ListenAndServe(addr, http.HandlerFunc(wrapCallback(cb)))
}

func wrapCallback(cb func(*Results, *http.Request) error) func(http.ResponseWriter, *http.Request) {
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
		sort.Sort(res.Answers)
		err := cb(res, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// Len implements the sort interface
func (s ResultsAnswerSlice) Len() int {
	return len(s)
}

// Less implements the sort interface
func (s ResultsAnswerSlice) Less(i, j int) bool {
	return s[i].FieldID < s[j].FieldID
}

// Swap implements the sort interface
func (s ResultsAnswerSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (a *ResultsAnswer) emptyValue(forJSON bool) interface{} {
	switch a.Type {
	case "number":
		if forJSON {
			return &jsonNumberValue{}
		}
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

func numberToInt64(n json.Number) (int64, error) {
	// attempt to use Int64, but if that doesn't work, fallback to Float64
	if i, err := n.Int64(); err == nil {
		return i, nil
	}
	f, err := n.Float64()
	if err != nil {
		return 0, err
	}
	return int64(f), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (a *ResultsAnswer) UnmarshalJSON(b []byte) error {
	ja := &jsonAnswer{}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	if err := d.Decode(ja); err != nil {
		return err
	}
	fid, err := numberToInt64(ja.jsonResultsAnswerMetadata.FieldID)
	if err != nil {
		return err
	}
	a.ResultsAnswerMetadata = ResultsAnswerMetadata{
		FieldID: fid,
		Type:    ja.jsonResultsAnswerMetadata.Type,
		Tags:    ja.jsonResultsAnswerMetadata.Tags,
	}

	a.Value = a.emptyValue(true)
	if a.Value != nil {
		d = json.NewDecoder(bytes.NewReader(ja.Value))
		d.UseNumber()
		if err := d.Decode(a.Value); err != nil {
			return err
		}
		switch v := a.Value.(type) {
		case *jsonNumberValue:
			nv, err := numberToInt64(v.Amount)
			if err != nil {
				return err
			}
			a.Value = &NumberValue{nv}
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

	a.Value = a.emptyValue(false)
	if a.Value != nil {
		if err := ba.Value.Unmarshal(a.Value); err != nil {
			return err
		}
	}
	return nil
}

// String returns a string version of the answer
func (a *ResultsAnswer) String() string {
	if a.Value == nil {
		return ""
	}
	switch v := a.Value.(type) {
	case *NumberValue:
		return strconv.FormatInt(v.Amount, 10)
	case *BooleanValue:
		if bool(*v) {
			return "true"
		} else {
			return "false"
		}
	case *ChoiceValue:
		if !v.EmptyOther {
			return v.Other
		}
		return v.Label
	case *ChoicesValue:
		if !v.EmptyOther {
			return v.Other
		}
		return strings.Join(v.Labels, ",")
	case *TextValue:
		return string(*v)
	default:
		llog.Error(fmt.Sprintf("encountered unknown type in stringValue: %T", a))
	}
	return ""
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
