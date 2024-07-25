package session

import "encoding/gob"

// RegisterSessionTypes registers those custom types
// that will be stored in the session.
func RegisterSessionTypes(types ...any) {
	for _, t := range types {
		gob.Register(t)
	}
}
