package form

import (
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
)

// use a single instance of Decoder, it caches struct info
var (
	// Shared decoder instance with default options from
	// the underlying library.
	decoder = form.NewDecoder()
)

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
