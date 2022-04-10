package pgdutil

import (
	"encoding/base64"
	"reflect"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

var (
	dateFormat = "2006-01-02"

	WrapCustomValidatorFunc = map[string]func(fl validator.FieldLevel) bool{
		"isRequiredWith": IsRequiredWith,
		"dateString":     DateString,
		"base64":         Base64,
	}
)

func IsRequiredWith(fl validator.FieldLevel) bool {
	field := fl.Field()
	otherField, _, _, _ := fl.GetStructFieldOK2()

	if otherField.IsValid() && otherField.Interface() != reflect.Zero(otherField.Type()).Interface() {
		if field.IsValid() && field.Interface() == reflect.Zero(field.Type()).Interface() {
			return false
		}
	}

	return true
}

func DateString(fl validator.FieldLevel) bool {
	field := fl.Field()
	vValue := fl.Param()

	if vValue == "" {
		vValue = dateFormat
	}

	if field.Interface() == reflect.Zero(field.Type()).Interface() {
		return true
	}

	date, err := time.Parse(vValue, field.Interface().(string))

	if err != nil || (date == time.Time{}) {
		return false
	}

	return true
}

func Base64(fl validator.FieldLevel) bool {
	field := fl.Field()

	// if field is nil
	if field.Interface() == reflect.Zero(field.Type()).Interface() {
		return true
	}

	// check field base64 or not, if error then false
	_, err := base64.StdEncoding.DecodeString(field.Interface().(string))

	return err == nil
}
