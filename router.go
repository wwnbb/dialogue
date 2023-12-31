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
			valid := validateAndExtractParams(d, pattern)
			if valid {
				return route(d)
			}
		}
		return d
	}
}

func validateAndExtractParams(d *Dialogue, pattern string) bool {
	path := d.Request.URL.Path
	partsPattern := strings.Split(pattern, "/")
	partsPath := strings.Split(path, "/")

	if len(partsPattern) > len(partsPath) {
		return false
	}

	params := map[string]Param{}

	for i, part := range partsPattern {
		if strings.HasPrefix(part, "<") && strings.HasSuffix(part, ">") {
			paramInfo := strings.Trim(part, "<>")
			paramParts := strings.Split(paramInfo, ":")
			paramName, paramType := paramParts[0], paramParts[1]

			var value interface{}
			var err error

			switch paramType {
			case "uuid4":
				value, err = uuid.Parse(partsPath[i])
			case "datetime":
				value, err = time.Parse(time.RFC3339, partsPath[i])
			case "int":
				value, err = strconv.Atoi(partsPath[i])
			case "float":
				value, err = strconv.ParseFloat(partsPath[i], 64)
			case "bool":
				value, err = strconv.ParseBool(partsPath[i])
			case "string":
				value, err = partsPath[i], nil // Strings are always valid
			case "hex":
				value, err = strconv.ParseUint(partsPath[i], 16, 64)
			default:
				log.Printf("Warning: Unknown type %v", paramType)
				return false
			}
			if err != nil {
				return false
			}

			params[paramName] = Param{Type: paramType, Value: value}
		} else if part != partsPath[i] {
			return false
		}
	}

	d.PathParams = params // Set the validated and extracted parameters
	return true
}
