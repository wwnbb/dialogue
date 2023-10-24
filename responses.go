package dialogue

import (
	"encoding/json"
	"net/http"
)

func SetStatusCode(d *Dialogue, statusCode int) *Dialogue {
	d.ResponseWriter.WriteHeader(statusCode)
	return d
}

func SetHeader(d *Dialogue, key, value string) *Dialogue {
	d.ResponseWriter.Header().Set(key, value)
	return d
}

func RedirectTo(d *Dialogue, url string, statusCode int) *Dialogue {
	http.Redirect(d.ResponseWriter, d.Request, url, statusCode)
	d.isProcessed = true
	return d
}

func SetContentType(d *Dialogue, contentType string) *Dialogue {
	return SetHeader(d, "Content-Type", contentType)
}

func WriteResponse(d *Dialogue, statusCode int, content []byte) *Dialogue {
	d = SetStatusCode(d, statusCode)
	_, _ = d.ResponseWriter.Write(content)
	d.isProcessed = true
	return d
}

func WriteResponseString(d *Dialogue, statusCode int, content string) *Dialogue {
	d = SetStatusCode(d, statusCode)
	_, _ = d.ResponseWriter.Write([]byte(content))
	d.isProcessed = true
	return d
}

func WriteResponseJson(d *Dialogue, statusCode int, content interface{}) *Dialogue {
	d = SetContentType(d, "application/json")
	jsonContent, err := json.Marshal(content)
	if err != nil {
		http.Error(d.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return d
	}
	return WriteResponse(d, statusCode, jsonContent)
}

func ServeFile(d *Dialogue, filePath string) *Dialogue {
	d = SetContentType(d, "application/octet-stream")
	http.ServeFile(d.ResponseWriter, d.Request, filePath)
	d.isProcessed = true
	return d
}

func WriteErrorResponse(d *Dialogue, statusCode int, errorMessage string) *Dialogue {
	d.ResponseWriter.WriteHeader(statusCode)
	_, _ = d.ResponseWriter.Write([]byte(errorMessage))
	d.isProcessed = true
	return d
}
