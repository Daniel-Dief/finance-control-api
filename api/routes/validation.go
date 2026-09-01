package routes

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

var validate = validator.New()

// bindAndValidate binds the request body into the destination struct and
// validates it according to the struct's `validate` tags. On failure it writes a
// 400 response and returns false so the caller can bail out.
func bindAndValidate(c *echo.Context, dst any) bool {
	if err := c.Bind(dst); err != nil {
		c.JSON(400, map[string]interface{}{"error": "Erro ao processar o corpo da requisição."})
		return false
	}

	if err := validate.Struct(dst); err != nil {
		c.JSON(400, map[string]interface{}{"error": validationErrorMessages(dst, err)})
		return false
	}

	return true
}

// bindAndValidateQuery binds the URL query parameters into the destination
// struct (tagged with `query`) and validates it. On failure it writes a 400
// response and returns false so the caller can bail out.
func bindAndValidateQuery(c *echo.Context, dst any) bool {
	if err := bindQuery(c, dst); err != nil {
		c.JSON(400, map[string]interface{}{"error": errorMessage(err)})
		return false
	}

	if err := validate.Struct(dst); err != nil {
		c.JSON(400, map[string]interface{}{"error": validationErrorMessages(dst, err)})
		return false
	}

	return true
}

// errorMessage returns a clean message for the given error. It extracts the
// Message from an *echo.HTTPError (ignoring the Code part) and falls back to
// err.Error() otherwise.
func errorMessage(err error) string {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return fmt.Sprintf("%v", httpErr.Message)
	}
	return err.Error()
}

// PathIDParams defines the `:id` path parameter used by the resource routes.
type PathIDParams struct {
	ID int `param:"id" validate:"required,gt=0"`
}

// bindAndValidatePath binds the URL path parameters into the destination
// struct (tagged with `param`) and validates it. On failure it writes a 400
// response and returns false so the caller can bail out.
func bindAndValidatePath(c *echo.Context, dst any) bool {
	if err := bindPath(c, dst); err != nil {
		c.JSON(400, map[string]interface{}{"error": "ID inválido"})
		return false
	}

	if err := validate.Struct(dst); err != nil {
		c.JSON(400, map[string]interface{}{"error": "ID inválido"})
		return false
	}

	return true
}

// bindPath maps URL path parameters into dst fields tagged with `param`.
func bindPath(c *echo.Context, dst any) error {
	v := reflect.ValueOf(dst).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("param")
		if tag == "" {
			continue
		}

		value := c.Param(tag)
		if value == "" {
			continue
		}

		target := v.Field(i)

		if target.Kind() == reflect.Pointer {
			target = reflect.New(target.Type().Elem()).Elem()
		}

		switch target.Kind() {
		case reflect.String:
			target.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return &echo.HTTPError{Message: "ID inválido"}
			}
			target.SetInt(num)
		}

		if v.Field(i).Kind() == reflect.Pointer {
			v.Field(i).Set(target.Addr())
		}
	}

	return nil
}

// bindQuery maps URL query parameters into dst fields tagged with `query`.
func bindQuery(c *echo.Context, dst any) error {
	v := reflect.ValueOf(dst).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("query")
		if tag == "" {
			continue
		}

		value := c.QueryParam(tag)
		if value == "" {
			continue
		}

		target := v.Field(i)

		// Pointer-typed fields stay nil when the parameter is absent and are
		// allocated on the fly when it is present.
		if target.Kind() == reflect.Pointer {
			target = reflect.New(target.Type().Elem()).Elem()
		}

		switch target.Kind() {
		case reflect.String:
			target.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return &echo.HTTPError{Message: "O parâmetro de consulta '" + tag + "' deve ser um número válido"}
			}
			target.SetInt(num)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			num, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return &echo.HTTPError{Message: "O parâmetro de consulta '" + tag + "' deve ser um número válido"}
			}
			target.SetUint(num)
		}

		if v.Field(i).Kind() == reflect.Pointer {
			v.Field(i).Set(target.Addr())
		}
	}

	return nil
}

// validationErrorMessages converts validator.ValidationErrors into a human
// readable slice of messages. The target struct is used to resolve friendly
// field labels from the `json` tag.
func validationErrorMessages(dst any, err error) []string {
	t := reflect.TypeOf(dst)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fieldLabel := func(fieldName string) string {
		if f, ok := t.FieldByName(fieldName); ok {
			if jsonTag := f.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				return jsonTag
			}
		}
		return fieldName
	}

	messages := make([]string, 0)
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return append(messages, err.Error())
	}

	for _, fe := range validationErrors {
		label := fieldLabel(fe.Field())
		switch fe.Tag() {
		case "required":
			messages = append(messages, "O campo '"+label+"' é obrigatório")
		case "email":
			messages = append(messages, "O campo '"+label+"' deve ser um e-mail válido")
		case "oneof":
			messages = append(messages, "O campo '"+label+"' deve ser um dos valores: "+fe.Param())
		case "datetime":
			messages = append(messages, "O campo '"+label+"' deve estar no formato "+fe.Param())
		case "min", "max":
			operator := "maior"
			if fe.Tag() == "max" {
				operator = "menor"
			}
			messages = append(messages, "O campo '"+label+"' deve ser "+operator+" ou igual a "+fe.Param())
		default:
			messages = append(messages, "O campo '"+label+"' é inválido")
		}
	}

	return messages
}
