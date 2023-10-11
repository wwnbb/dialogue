package dialogue

import (
	"log"
)

type RouteMap map[string]DialogueFunc

type DialogueFunc func(*Dialogue) *Dialogue

// Chain takes a list of DialogueFunc and returns a DialogueFunc.
func Chain(funcVariadic ...DialogueFunc) DialogueFunc {
	return func(d *Dialogue) *Dialogue {
		for _, f := range funcVariadic {
			d = d.Map(f)
		}
		return d
	}
}

// Map is a method that takes a DialogueFunc and returns a new Dialogue.
// It is used to chain middleware and handle requests.
func (d *Dialogue) Map(f DialogueFunc) *Dialogue {
	if d.isProcessed {
		return d
	}
	return f(d)
}

func (d Dialogue) Finish() {
	if d.isProcessed {
		log.Printf("Warning: Request %v was processed before", d.Request.URL.Path)
	}
	d.isProcessed = true
}

func Switch(routes RouteMap) DialogueFunc {
	return func(d *Dialogue) *Dialogue {
		if handler, exists := routes[d.Request.URL.Path]; exists {
			return handler(d)
		}
		return d
	}
}
