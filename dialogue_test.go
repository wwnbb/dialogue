package dialogue

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFoundHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/non-existent", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	handler := NotFoundHandler()

	d := NewDialogue(req, rec)
	handler(d)

	if status := rec.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "404 page not found\n"
	if rec.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}

	if !d.isProcessed {
		t.Errorf("Dialogue was not marked as processed")
	}

	if len(d.Logs) == 0 || d.Logs[0] != "Sent 404 Not Found" {
		t.Errorf("Logging did not work as expected")
	}
}

func TestListenAndServe(t *testing.T) {
	addr := ":8080"
	handler := NotFoundHandler() // Or any other handler you want to test

	server, err := ListenAndServe(addr, handler)
	if err != nil {
		t.Fatalf("Could not start server: %v", err)
	}

	defer server.Close()

	// Now you can send requests to the server and check the responses
}
