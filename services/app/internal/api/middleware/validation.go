package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
)

func Validation[T any]() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data T
		var err error

		if ctx.ContentType() == gin.MIMEJSON {
			err = ctx.ShouldBindBodyWith(&data, binding.JSON)
		} else {
			err = ctx.ShouldBind(&data)
		}

		if err != nil {
			message, details := humanizeValidationError(err)

			ctx.AbortWithStatusJSON(http.StatusBadRequest, &shared.HttpError{
				Code:    http.StatusBadRequest,
				Message: message,
				Details: details,
			})
			return
		}

		ctx.Set("validatedBody", data)
		ctx.Next()
	}
}

func humanizeValidationError(err error) (string, map[string]any) {
	var validationErrs validator.ValidationErrors

	if errors.As(err, &validationErrs) {
		errorsList := make([]string, 0, len(validationErrs))

		for _, validationErr := range validationErrs {
			errorsList = append(errorsList, validatorTagMessage(validationErr))
		}

		return "Validation failed", map[string]any{
			"errors": errorsList,
		}
	}

	return "Invalid request body", map[string]any{
		"errors": []string{err.Error()},
	}
}

func validatorTagMessage(validationErr validator.FieldError) string {
	fieldName := humanizeField(validationErr.Field())

	switch validationErr.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldName)
	case "eqfield":
		otherField := humanizeField(validationErr.Param())
		return fmt.Sprintf("%s must match %s", fieldName, otherField)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fieldName)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fieldName)
	default:
		return fmt.Sprintf("%s is invalid", fieldName)
	}
}

func humanizeField(fieldName string) string {
	if fieldName == "" {
		return "field"
	}

	b := strings.Builder{}

	for i, r := range fieldName {
		if i > 0 && unicode.IsUpper(r) {
			b.WriteRune(' ')
		}

		b.WriteRune(unicode.ToLower(r))
	}

	return b.String()
}
