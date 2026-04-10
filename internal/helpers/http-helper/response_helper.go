package http_helper

import (
	"io"
	"net/http"
	"testing"

	"api-tests-template/internal/utils"

	"github.com/stretchr/testify/require"
)

// ReadResponseBody читает body из response и закрывает connection
func ReadResponseBody(t *testing.T, resp *http.Response) string {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Ошибка чтения body из ответа")

	return string(body)
}

// AssertStatusCode проверяет status code и возвращает body
func AssertStatusCode(t *testing.T, resp *http.Response, expectedStatusCode int) string {
	body := ReadResponseBody(t, resp)

	require.Equalf(t, expectedStatusCode, resp.StatusCode,
		"Ожидался HTTP status %d, получен %d. Body: %s",
		expectedStatusCode, resp.StatusCode, body)

	utils.LogWithLabelAndTimeStamp("HTTP", "Status: "+resp.Status)

	return body
}
