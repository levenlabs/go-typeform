package fields

import (
	"github.com/levenlabs/golib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/validator.v2"
	. "testing"
)

func randField(t FieldType) Field {
	return Field{
		Type:     t,
		Question: testutil.RandStr(),
	}
}

func randChoices(l int) []MultipleChoiceChoice {
	d := make([]MultipleChoiceChoice, l)
	for i := range d {
		d[i] = MultipleChoiceChoice{testutil.RandStr()}
	}
	return d
}

func TestMultipleChoiceChoice(t *T) {
	// cannot have more than 512 characters
	assert.NotNil(t, validator.Validate(&MultipleChoiceChoice{
		string(make([]byte, 513)),
	}))

	// cannot be empty
	assert.NotNil(t, validator.Validate(&MultipleChoiceChoice{
		"",
	}))

	assert.Nil(t, validator.Validate(&MultipleChoiceChoice{
		testutil.RandStr(),
	}))
}

func TestMultipleChoice(t *T) {
	// there is at least 1 choice required
	assert.NotNil(t, validator.Validate(&MultipleChoice{
		Field:   randField(TypeStatement),
		Choices: []MultipleChoiceChoice{},
	}))

	// you cannot have more than 25 choices
	assert.NotNil(t, validator.Validate(&MultipleChoice{
		Field:   randField(TypeStatement),
		Choices: randChoices(26),
	}))

	assert.Nil(t, validator.Validate(&MultipleChoice{
		Field:   randField(TypeStatement),
		Choices: randChoices(1),
	}))
}

func TestOpinionLabels(t *T) {
	// cannot have a left of > 100
	assert.NotNil(t, validator.Validate(&OpinionLabels{
		Left: string(make([]byte, 101)),
	}))

	// cannot have a center of > 100
	assert.NotNil(t, validator.Validate(&OpinionLabels{
		Center: string(make([]byte, 101)),
	}))

	// cannot have a right of > 100
	assert.NotNil(t, validator.Validate(&OpinionLabels{
		Right: string(make([]byte, 101)),
	}))
}

func TestOpinionScale(t *T) {
	// steps must be >= 5
	assert.NotNil(t, validator.Validate(&OpinionScale{
		Field: randField(TypeOpinionScale),
	}))

	// steps must be <= 11
	assert.NotNil(t, validator.Validate(&OpinionScale{
		Field: randField(TypeOpinionScale),
		Steps: 12,
	}))

	assert.Nil(t, validator.Validate(&OpinionScale{
		Field: randField(TypeStatement),
		Steps: 5,
	}))
}

func TestFieldInterface(t *T) {
	s := interface{}(&Statement{})
	_, ok := s.(FieldInterface)
	assert.True(t, ok)

	mc := interface{}(&MultipleChoice{})
	_, ok = mc.(FieldInterface)
	assert.True(t, ok)

	os := interface{}(&OpinionScale{})
	_, ok = os.(FieldInterface)
	assert.True(t, ok)

	f := &Field{
		Type:        TypeStatement,
		Question:    "Hey?",
		Ref:         "hey",
		Description: "This is a test",
		Required:    true,
	}
	fi, ok := interface{}(f).(FieldInterface)
	require.True(t, ok)

	assert.Equal(t, TypeStatement, fi.GetType())
	assert.Equal(t, "Hey?", fi.GetQuestion())
	assert.Equal(t, "hey", fi.GetRef())
	assert.Equal(t, "This is a test", fi.GetDescription())
	assert.Equal(t, true, fi.GetRequired())
	fi.SetValue("val")
	assert.Equal(t, "val", f.Value)
	assert.Equal(t, "val", f.GetValue())
}
