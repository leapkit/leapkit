package form

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/gofrs/uuid/v5"
)

// use a single instance of Decoder, it caches struct info
var (
	// Shared decoder instance with default options from
	// the underlying library.
	decoder = form.NewDecoder()
)

func init() {
	// Register custom and common type decoder
	// functions.
	decoder.RegisterCustomTypeFunc(decodeUUID, uuid.UUID{})
	decoder.RegisterCustomTypeFunc(decodeUUIDSlice, []uuid.UUID{})
}

// RegisterCustomTypeFunc registers a custom type decoder func for a type.
// This is useful when you want to use a custom type or a type from an external
// package like uuid.UUID and want to decode it from a string.
func RegisterCustomTypeFunc(fn form.DecodeCustomTypeFunc, kind interface{}) {
	decoder.RegisterCustomTypeFunc(fn, kind)
}

// Decode decodes the request body into dst, which must be a pointer of a struct.
// If there is no body or the body is empty, it will take the query string as the
// body. If the Content-Type is multipart/form-data.
func Decode(r *http.Request, dst interface{}) error {
	//MultipartForm
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			return err
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			return err
		}
	}

	data := r.Form
	if len(data) == 0 {
		r.Form = r.URL.Query()
	}

	err := decoder.Decode(dst, r.Form)
	return err
}

// decodeUUID a single uuid from a string
// and returns an error if there is a problem
func decodeUUID(vals []string) (interface{}, error) {
	uu, err := uuid.FromString(vals[0])
	if err != nil {
		err = fmt.Errorf("error parsing uuid: %w", err)
	}

	return uu, err
}

// decodeUUIDSlice decodes a slice of uuids from a string
// and returns an error if there is a problem
func decodeUUIDSlice(vals []string) (interface{}, error) {
	var uus []uuid.UUID

	for _, val := range vals {
		uuid, err := uuid.FromString(val)
		if err != nil {
			err = fmt.Errorf("error parsing uuid: %w", err)
			return nil, err
		}

		uus = append(uus, uuid)
	}

	return uus, nil
}
