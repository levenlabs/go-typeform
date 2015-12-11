# go-typeform

Libraries for interacting with [Typeform I/O](http://docs.typeform.io/docs).

## fields

This package contains all the fields needed to represent a form and different
fields that you send/receive from the API. A form could be many different types
of fields and each field has different properties. The `fields` package exports
a `Form` struct that can be JSON or BSON marshal'd and unmarshal'd while keepin
all the field-specific data for each type. There are also validations using the
[validator](https://github.com/go-validator/validator) library for each of the
properties of each specific field.

**Not all fields are implemented yet. This is a WIP**