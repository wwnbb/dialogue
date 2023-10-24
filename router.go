package dialogue

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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
		for pattern, route := range routes {
			valid := validateAndExtractParams(pattern, d.Request.URL.Path)
			if valid {
				return route(d)
			}
		}
		return d
	}
}

func validateAndExtractParams(pattern string, path string) bool {
	partsPattern := strings.Split(pattern, "/")
	partsPath := strings.Split(path, "/")

	if len(partsPattern) > len(partsPath) {
		return false
	}

	params := map[string]string{}

	for i, part := range partsPattern {
		if strings.HasPrefix(part, "<") && strings.HasSuffix(part, ">") {
			paramInfo := strings.Trim(part, "<>")
			paramParts := strings.Split(paramInfo, ":")
			paramName, paramType := paramParts[0], paramParts[1]

			switch paramType {
			case "uuid4":
				if _, err := uuid.Parse(partsPath[i]); err != nil {
					return false
				}
			case "datetime":
				if _, err := time.Parse(time.RFC3339, partsPath[i]); err != nil {
					return false
				}
			case "int":
				if _, err := strconv.Atoi(partsPath[i]); err != nil {
					return false
				}
			case "float":
				if _, err := strconv.ParseFloat(partsPath[i], 64); err != nil {
					return false
				}
			case "bool":
				if _, err := strconv.ParseBool(partsPath[i]); err != nil {
					return false
				}
			case "string":
				// No validation needed for string
			case "hex":
				if _, err := strconv.ParseUint(partsPath[i], 16, 64); err != nil {
					return false
				}
			default:
				return false
			}

			params[paramName] = partsPath[i]
		} else if part != partsPath[i] {
			return false
		}
	}

	return true
}
