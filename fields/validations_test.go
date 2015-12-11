package fields

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/validator.v2"
	. "testing"
)

func TestURL(t *T) {
	tags := "validateURL"
	assert.Nil(t, validator.Valid("http://typeform.io/", tags))
	assert.NotNil(t, validator.Valid("11111", tags))
	assert.NotNil(t, validator.Valid(11111, tags))
	assert.Nil(t, validator.Valid("", tags))
}
