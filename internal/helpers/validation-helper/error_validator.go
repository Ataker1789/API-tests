package validation_helper

import (
	"testing"

	"api-tests-template/internal/utils"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// ValidateErrorResponse проверяет наличие полей error и message в ответе
func ValidateErrorResponse(t *testing.T, body string) (string, string) {
	errorField := gjson.Get(body, "error").String()
	messageField := gjson.Get(body, "message").String()

	require.NotEmpty(t, errorField, "Поле 'error' должно быть заполнено")
	require.NotEmpty(t, messageField, "Поле 'message' должно быть заполнено")

	utils.LogWithLabelAndTimeStamp("Validation", "Error: "+errorField+", Message: "+messageField)

	return errorField, messageField
}

// ValidateUnauthorizedError проверяет ошибку unauthorized (401)
func ValidateUnauthorizedError(t *testing.T, body string) {
	errorField, _ := ValidateErrorResponse(t, body)
	require.Equal(t, "unauthorized", errorField, "Должна быть ошибка 'unauthorized'")
}

// ValidateBadRequestError проверяет ошибку bad_request (400)
func ValidateBadRequestError(t *testing.T, body string) {
	errorField, messageField := ValidateErrorResponse(t, body)
	require.NotEmpty(t, errorField, "Поле 'error' не должно быть пустым при 400")
	require.NotEmpty(t, messageField, "Поле 'message' не должно быть пустым при 400")
}
