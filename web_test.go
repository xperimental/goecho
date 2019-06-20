package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchoHandler(t *testing.T) {
	for _, test := range []struct {
		hostname string
		env      []string
		url      string
		header   http.Header
		status   int
		body     string
	}{
		{
			url:    "/",
			status: http.StatusOK,
			body:   "URL: /\nHeader:\n",
		},
		{
			url:    "/",
			header: http.Header{},
			status: http.StatusOK,
			body:   "URL: /\nHeader:\n",
		},
		{
			url: "/",
			header: http.Header{
				"key1": []string{"value1"},
			},
			status: http.StatusOK,
			body:   "URL: /\nHeader:\nkey1 -> \"value1\"\n",
		},
		{
			url: "/",
			header: http.Header{
				"key1": []string{"value1", "value2"},
			},
			status: http.StatusOK,
			body:   "URL: /\nHeader:\nkey1 -> \"value1\"; \"value2\"\n",
		},
		{
			url: "/",
			header: http.Header{
				"key1": []string{"value1"},
				"key2": []string{"value2"},
			},
			status: http.StatusOK,
			body:   "URL: /\nHeader:\nkey1 -> \"value1\"\nkey2 -> \"value2\"\n",
		},
		{
			url: "/",
			header: http.Header{
				"key2": []string{"value2"},
				"key1": []string{"value1"},
			},
			status: http.StatusOK,
			body:   "URL: /\nHeader:\nkey1 -> \"value1\"\nkey2 -> \"value2\"\n",
		},
		{
			hostname: "test",
			url:      "/",
			status:   http.StatusOK,
			body:     "URL: /\nHeader:\n\nServer: test\n",
		},
		{
			env: []string{
				"key=value",
				"key2=value2",
			},
			url:    "/",
			status: http.StatusOK,
			body:   "URL: /\nHeader:\n",
		},
		{
			env: []string{
				"key=value",
				"key2=value2",
			},
			url:    "/?env=true",
			status: http.StatusOK,
			body:   "URL: /?env=true\nHeader:\n\nEnvironment:\nkey=value\nkey2=value2\n",
		},
	} {
		req, err := http.NewRequest(http.MethodGet, test.url, nil)
		if err != nil {
			t.Fatalf("can not create request: %s", err)
		}
		req.Header = test.header

		w := httptest.NewRecorder()

		echoHandler(test.hostname, test.env).ServeHTTP(w, req)

		if w.Code != test.status {
			t.Errorf("got status %d, expected %d", w.Code, test.status)
		}

		body := w.Body.String()
		if body != test.body {
			t.Errorf("got '%s', wanted '%s'", body, test.body)
		}
	}
}

func TestVersionHandler(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedContent := "application/json; chatset=utf-8"

	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatalf("can not create request: %s", err)
	}

	w := httptest.NewRecorder()
	handler := versionHandler("test-version")

	handler.ServeHTTP(w, req)

	if w.Code != expectedStatus {
		t.Errorf("got status %d, wanted %d", w.Code, expectedStatus)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != expectedContent {
		t.Errorf("got type '%s', wanted '%s'", contentType, expectedContent)
	}

	body := w.Body.String()
	if len(body) == 0 {
		t.Error("got no body.")
	}
}

func TestReadyHandler(t *testing.T) {
	readyness, unready := readyHandler()

	req := httptest.NewRequest(http.MethodGet, "/_ready", nil)
	rec := httptest.NewRecorder()

	readyness.ServeHTTP(rec, req)

	expectedStatus := http.StatusOK
	if rec.Code != expectedStatus {
		t.Errorf("got status %d, wanted %d", rec.Code, expectedStatus)
	}

	unready()

	rec = httptest.NewRecorder()
	readyness.ServeHTTP(rec, req)

	expectedStatus = http.StatusInternalServerError
	if rec.Code != expectedStatus {
		t.Errorf("got status %d, wanted %d", rec.Code, expectedStatus)
	}
}

func TestCreateServer(t *testing.T) {
	addr := "localhost:3000"
	s, _ := createServer(addr, "", "", []string{})

	if s.Addr != addr {
		t.Errorf("got '%s', wanted '%s'", s.Addr, addr)
	}

	if s.Handler == nil {
		t.Errorf("handler not set")
	}
}
