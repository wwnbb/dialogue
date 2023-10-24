package dialogue

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetStatusCode(t *testing.T) {
	rec := httptest.NewRecorder()
	d := &Dialogue{ResponseWriter: rec}
	d = SetStatusCode(d, http.StatusTeapot)
	if rec.Code != http.StatusTeapot {
		t.Errorf("expected status code %v, got %v", http.StatusTeapot, rec.Code)
	}
}

func TestSetHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	d := &Dialogue{ResponseWriter: rec}
	d = SetHeader(d, "X-Test-Header", "test-value")
	if rec.Header().Get("X-Test-Header") != "test-value" {
		t.Errorf("expected header value 'test-value', got %v", rec.Header().Get("X-Test-Header"))
	}
}

func TestWriteResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	d := &Dialogue{ResponseWriter: rec}
	content := []byte("hello world")
	d = WriteResponse(d, http.StatusOK, content)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != string(content) {
		t.Errorf("expected body '%v', got '%v'", string(content), rec.Body.String())
	}
}

func TestWriteResponseString(t *testing.T) {
	rec := httptest.NewRecorder()
	d := &Dialogue{ResponseWriter: rec}
	content := "hello world"
	d = WriteResponseString(d, http.StatusOK, content)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != content {
		t.Errorf("expected body '%v', got '%v'", content, rec.Body.String())
	}
}

func TestWriteResponseJson(t *testing.T) {
	rec := httptest.NewRecorder()
	d := &Dialogue{ResponseWriter: rec}
	content := map[string]string{"message": "hello world"}
	d = WriteResponseJson(d, http.StatusOK, content)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, rec.Code)
	}
	expectedJson := `{"message":"hello world"}`
	if rec.Body.String() != expectedJson {
		t.Errorf("expected body '%v', got '%v'", expectedJson, rec.Body.String())
	}
}

func TestServeFile(t *testing.T) {
	rec := httptest.NewRecorder()
	d := &Dialogue{ResponseWriter: rec, Request: httptest.NewRequest(http.MethodGet, "/file.txt", nil)}
	filePath := "test_resources/example_file.txt" // Assume this file exists with content "hello world"
	d = ServeFile(d, filePath)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "hello world\n" { // Assuming file.txt contains "hello world"
		t.Errorf("expected body 'hello world', got '%v'", rec.Body.String())
	}
}

func TestWriteErrorResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	d := &Dialogue{ResponseWriter: rec}
	errorMessage := "error occurred"
	d = WriteErrorResponse(d, http.StatusInternalServerError, errorMessage)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %v, got %v", http.StatusInternalServerError, rec.Code)
	}
	if rec.Body.String() != errorMessage {
		t.Errorf("expected body '%v', got '%v'", errorMessage, rec.Body.String())
	}
}

func TestRedirectTo(t *testing.T) {
	rec := httptest.NewRecorder()
	d := &Dialogue{ResponseWriter: rec, Request: httptest.NewRequest(http.MethodGet, "/", nil)}
	d = RedirectTo(d, "/redirect", http.StatusMovedPermanently)
	if rec.Code != http.StatusMovedPermanently {
		t.Errorf("expected status code %v, got %v", http.StatusMovedPermanently, rec.Code)
	}
	location, _ := rec.Result().Location()
	if location.String() != "/redirect" {
		t.Errorf("expected location '/redirect', got '%v'", location.String())
	}
}
