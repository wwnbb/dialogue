// Dialogue is a minimalistic web framework for Go. Based on idea of eliminating global state,
// and using monadic pattern to handle requests/response data flow, logging connections errors and database connections.
package dialogue

import (
	"log"
	"net/http"
	"time"
)

type Dialogue struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Logs           []string
	isProcessed    bool
}

func NewDialogue(r *http.Request, w http.ResponseWriter) *Dialogue {
	return &Dialogue{
		Request:        r,
		ResponseWriter: w,
		Logs:           make([]string, 0),
	}
}

func NotFoundHandler() DialogueFunc {
	return func(d *Dialogue) *Dialogue {
		http.NotFound(d.ResponseWriter, d.Request)
		d.isProcessed = true // marking as processed since a response is given
		d.Logs = append(d.Logs, "Sent 404 Not Found")
		return d
	}
}

func ListenAndServe(addr string, dialogueHandler DialogueFunc) (*http.Server, error) {
	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := NewDialogue(r, w)
		dialogueHandler(d)
	})

	s := &http.Server{
		Addr:           addr,
		Handler:        httpHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	return s, nil
}
