package pgdutil

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type (
	DummyEcho struct {
		EchoObj  *echo.Echo
		Request  *http.Request
		Response *httptest.ResponseRecorder
		Context  echo.Context
	}

	customValidator struct {
		Validator *validator.Validate
	}
)

func NewDummyEcho(method, path string, pl ...interface{}) DummyEcho {
	var body string
	e := echo.New()
	cv := customValidator{}
	cv.CustomValidation()
	e.Validator = &cv

	if pl != nil {
		body = ToJson(pl[0])
	}

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)

	return DummyEcho{
		EchoObj:  e,
		Request:  req,
		Response: resp,
		Context:  c,
	}
}

func (cv *customValidator) CustomValidation() {
	validator := validator.New()

	for key, fn := range WrapCustomValidatorFunc {
		_ = validator.RegisterValidation(key, fn)
	}

	cv.Validator = validator
}

func (cv *customValidator) Validate(i interface{}) error {
	rValue := reflect.ValueOf(i)
	if rValue.Kind() == reflect.Ptr {
		rValue = rValue.Elem()
	}

	if rValue.Kind() == reflect.Struct {
		return cv.Validator.Struct(i)
	}

	obj, ok := i.(map[string]interface{})

	if !ok {
		return ErrInternalServerError
	}

	isError, ok := obj["isError"].(bool)

	if !ok {
		return nil
	}

	if isError {
		return ErrInternalServerError
	}

	return nil
}
