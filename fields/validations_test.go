package fields

import (
	. "testing"
	"github.com/stretchr/testify/assert"
	"gopkg.in/validator.v2"
)

func TestURL(t *T) {
	tags := "validateURL"
	assert.Nil(t, validator.Valid("http://typeform.io/", tags))
	assert.NotNil(t, validator.Valid("11111", tags))
	assert.NotNil(t, validator.Valid(11111, tags))
}
