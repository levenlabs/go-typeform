package fields

import (
	"errors"
	"github.com/levenlabs/golib/rpcutil"
	"gopkg.in/validator.v2"
	"net/url"
)

func init() {
	rpcutil.InstallCustomValidators()

	// validateURL validates that the field is a url
	validator.SetValidationFunc("validateURL", validateURL)
}

func validateURL(v interface{}, _ string) error {
	u, ok := v.(string)
	if !ok {
		return validator.ErrUnsupported
	}
	uo, err := url.Parse(u)
	if uo.Scheme == "" {
		return errors.New("missing scheme in url")
	}
	if uo.Host == "" {
		return errors.New("missing host in url")
	}
	return err
}
