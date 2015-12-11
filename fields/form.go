package fields

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

// FormMetadata represents everything about a Form except for the Fields
type FormMetadata struct {
	Title      string         `json:"title"                        bson:"t"           validate:"min=1,max=256"`
	Tags       []string       `json:"tags,omitempty"               bson:"g,omitempty" validate:"arrMap=min=1,arrMap=max=128,max=100"`
	WebhookURL string         `json:"webhook_submit_url,omitempty" bson:"w,omitempty" validate:"validateURL"`
}

// A Form is a group of Fields that can be submitted to TypeForm's [/forms](http://docs.typeform.io/docs/forms)
// endpoint
type Form struct {
	FormMetadata                                                  `bson:",inline"`
	Fields     []interface{}  `json:"fields"                       bson:"f"           validate:"min=1,max=500"`
}

// jsonForm is used to Unmarshal into since it has Fields of json.RawMessage
type jsonForm struct {
	FormMetadata
	Fields []json.RawMessage  `json:"fields"                       bson:"f"`
}

// bsonForm is used to Unmarshal into since it has Fields of bson.Raw
type bsonForm struct {
	FormMetadata
	Fields []bson.Raw         `json:"fields"                       bson:"f"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (f *Form) UnmarshalJSON(b []byte) error {
	jf := &jsonForm{}
	if err := json.Unmarshal(b, jf); err != nil {
		return err
	}
	f.FormMetadata = jf.FormMetadata

	var err error
	var dst interface{}
	f.Fields = make([]interface{}, len(jf.Fields))
	for i, qstr := range jf.Fields {
		q := &Field{}
		if err = json.Unmarshal(qstr, q); err != nil {
			break
		}
		dst = q.emptyInterface()
		if err = json.Unmarshal(qstr, dst); err != nil {
			break
		}
		f.Fields[i] = dst
	}
	return err
}

// SetBSON implements the bson.Setter interface
func (f *Form) SetBSON(raw bson.Raw) error {
	bf := &bsonForm{}
	if err := raw.Unmarshal(bf); err != nil {
		return err
	}
	f.FormMetadata = bf.FormMetadata

	var err error
	var dst interface{}
	f.Fields = make([]interface{}, len(bf.Fields))
	for i, qstr := range bf.Fields {
		q := &Field{}
		if err = qstr.Unmarshal(q); err != nil {
			break
		}
		dst = q.emptyInterface()
		if err = qstr.Unmarshal(dst); err != nil {
			break
		}
		f.Fields[i] = dst
	}
	return err
}

