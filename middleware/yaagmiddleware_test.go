package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/CoryARamirez/yaag/yaag/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/websocket"
	"github.com/CoryARamirez/yaag/yaag"
	"strings"
	"os"
)

func TestMain(m *testing.M)  {
	yaag.Init(&yaag.Config{
		On:       true,
		DocTitle: "Test",
		DocPath:  "apidoc.html",
		BaseUrls: map[string]string{"Production": "", "Staging": ""},
	})
	os.Exit(m.Run())
}

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
		next := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.code)
			w.Header().Set(test.contentTypeKey, test.contentTypeValue)
			w.Write(body.Bytes())
		}
		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		apiCall := &models.ApiCall{}
		outputRecorder := NewResponseRecorder(httptest.NewRecorder())
		Before(apiCall, request)
		next(outputRecorder, request)
		After(apiCall, outputRecorder, request)

		if outputRecorder.Status != test.code {
			t.Errorf("expected code to be %d, was %s", test.code, outputRecorder.Status)
		}

		contentType := outputRecorder.Header().Get(test.contentTypeKey)
		if contentType != test.contentTypeValue {
			t.Errorf("expected header %s to be %s, was %s", test.contentTypeKey, test.contentTypeValue, contentType)
		}

	}

}

func TestWithWebsocketUpgrade(t *testing.T) {
	upgrader := websocket.Upgrader{}
	next := func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			c.Close()
			t.Fatal(err, "upgrade err")
		}

		defer c.Close()
		mt, message, err := c.ReadMessage()
		t.Log(err, "readMessage")
		t.Log(c.WriteMessage(mt, message), "write err")

	}
	s := httptest.NewServer(HandleFunc(next))
	c, _, err := websocket.DefaultDialer.Dial(strings.Replace(
		s.URL, "http","ws",1,
	), nil)
	if err != nil {
		t.Fatal("dial:", err)
	}
	defer c.Close()
	err = c.WriteMessage(websocket.TextMessage, []byte("hello\n"))
	if err != nil {
		t.Log("write:", err)
		return
	}
	_, message, err := c.ReadMessage()
	if err != nil {
		t.Fatalf("read:", err)
		return
	}
	t.Log("recv: %s", message)
	if string(message) != "hello\n" {
		t.Fatal("Message doesn't match")
	}
}
