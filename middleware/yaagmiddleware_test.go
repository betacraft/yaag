package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/betacraft/yaag/yaag/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAfterSetContentType(t *testing.T) {

	type test struct {
		code             int
		contentTypeKey   string
		contentTypeValue string
	}

	tests := []test{
		{http.StatusOK, "Content-Type", "application/json"},
		{http.StatusInternalServerError, "Content-Type", "application/json"},
		{http.StatusBadRequest, "Content-Type", "application/json"},
		{http.StatusNotFound, "Content-Type", "application/json"},
	}

	testResponseBody := map[string]string{"test": "yo"}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(testResponseBody); err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		record := httptest.NewRecorder()
		record.Code = test.code
		record.Header().Set(test.contentTypeKey, test.contentTypeValue)
		record.Body = body

		outputRecorder := httptest.NewRecorder()
		After(&models.ApiCall{}, record, outputRecorder, request)

		if outputRecorder.Code != test.code {
			t.Errorf("expected code to be %d, was %s", test.code, outputRecorder.Code)
		}

		contentType := outputRecorder.Header().Get(test.contentTypeKey)
		if contentType != test.contentTypeValue {
			t.Errorf("expected header %s to be %s, was %s", test.contentTypeKey, test.contentTypeValue, contentType)
		}

	}

}
