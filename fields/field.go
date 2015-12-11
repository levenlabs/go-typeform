// The fields package contains all the fields to represent forms and their
// associated metadata.
// Each of these matches a field in the [Fields](http://docs.typeform.io/docs/introduction)
// section of their api documentation.
package fields

// Field is a generic Field that holds common properties of all Fields in a Form
type Field struct {
	Type        FieldType    `json:"type"                  bson:"t"`
	Question    string       `json:"question"              bson:"q"              validate:"nonzero,max=512"`
	Ref         string       `json:"ref,omitempty"         bson:"ref,omitempty"  validate:"max=128"`
	Description string       `json:"description,omitempty" bson:"d"              validate:"max=512"`
	Required    bool         `json:"required,omitempty"    bson:"r,omitempty"`
}

// OpinionLabels represents a OpinionScale's labels property
// It contains a label for the left, center, and right sides of the scale
type OpinionLabels struct {
	Left   string `json:"left,omitempty"                   bson:"l,omitempty"    validate:"max=100"`
	Center string `json:"center,omitempty"                 bson:"c,omitempty"    validate:"max=100"`
	Right  string `json:"right,omitempty"                  bson:"r,omitempty"    validate:"max=100"`
}

// OpinionScale is a scale from 0-(steps - 1)
type OpinionScale struct {
	Field                                                  `bson:",inline"`
	Steps      int64         `json:"steps"                  bson:"s"             validate:"min=5,max=11"`
	StartAtOne bool          `json:"start_at_one,omitempty" bson:"sao,omitempty"`
	// todo: figure out how to validate if completely empty (once https://github.com/go-validator/validator/pull/38)
	Labels     OpinionLabels `json:"labels,omitempty"       bson:"l,omitempty"`
}

// MultipleChoiceChoice is a choice in a MultipleChoice's Choices slice
type MultipleChoiceChoice struct {
	Label string `json:"label"                              bson:"l"             validate:"nonzero, max=512"`
}

// MultipleChoice is a question that contains multiple choices
type MultipleChoice struct {
	Field                                                  `bson:",inline"`
	Choices []MultipleChoiceChoice `json:"choices"          bson:"c"             validate:"min=1,max=25"`
}

// Statement is just text, no question to answer
type Statement struct {
	Field                                                  `bson:",inline"`
}

// FieldType describes the type of field
type FieldType string

var (
	TypeStatement FieldType = "statement"
	TypeOpinionScale FieldType = "opinion_scale"
	TypeMultipleChoice FieldType = "multiple_choice"
)

// emptyInterface can be used to get an empty specific struct for the type of
// field that the field is. This is used by the Form to unmarshal
func (q *Field) emptyInterface() interface{} {
	var dst interface{}
	switch q.Type {
	case TypeStatement:
		dst = &Statement{}
	case TypeOpinionScale:
		dst = &OpinionScale{}
	case TypeMultipleChoice:
		dst = &MultipleChoice{}
	default:
		dst = q
	}
	return dst
}
