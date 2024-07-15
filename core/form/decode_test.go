package form_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/leapkit/leapkit/core/form"
)

func TestRegisterCustomDecoder(t *testing.T) {
	vals := url.Values{
		"Ddd": []string{"21-01-01"},
		"Sss": []string{"hello"},
	}

	tr, err := http.NewRequest("POST", "/", strings.NewReader(vals.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	tr.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Registering custom type
	form.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.Parse("06-01-02", vals[0])
	}, time.Time{})

	st := struct {
		Ddd time.Time `form:"Ddd"`
		Sss string    `form:"Sss"`
	}{}

	err = form.Decode(tr, &st)
	if err != nil {
		t.Fatal(err)
	}

	if st.Ddd.Format("2006-01-02") != "2021-01-01" {
		t.Fatalf("expected 2021-01-01, got %v", st.Ddd.Format("2006-01-02"))
	}
}

func TestDecodeGet(t *testing.T) {
	vals := url.Values{
		"Ddd": []string{"21-01-01"},
		"Sss": []string{"hello"},
	}

	tr, err := http.NewRequest("GET", "/?"+vals.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	tr.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Registering custom type
	form.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.Parse("06-01-02", vals[0])
	}, time.Time{})

	st := struct {
		Ddd time.Time `form:"Ddd"`
		Sss string    `form:"Sss"`
	}{}

	err = form.Decode(tr, &st)
	if err != nil {
		t.Fatal(err)
	}

	if st.Ddd.Format("2006-01-02") != "2021-01-01" {
		t.Fatalf("expected 2021-01-01, got %v", st.Ddd.Format("2006-01-02"))
	}
}

func TestDecodeMultipartForm(t *testing.T) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	err := writer.WriteField("Ddd", "21-01-01")
	if err != nil {
		t.Fatal(err)
	}
	err = writer.WriteField("Sss", "hello")
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()

	tr, err := http.NewRequest("POST", "/", &buf)
	if err != nil {
		t.Fatal(err)
	}

	tr.Header.Set("Content-Type", writer.FormDataContentType())

	// Registering custom type
	form.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.Parse("06-01-02", vals[0])
	}, time.Time{})

	st := struct {
		Ddd time.Time `form:"Ddd"`
		Sss string    `form:"Sss"`
	}{}

	err = form.Decode(tr, &st)
	if err != nil {
		t.Fatal(err)
	}

	if st.Ddd.Format("2006-01-02") != "2021-01-01" {
		t.Fatalf("expected 2021-01-01, got %v", st.Ddd.Format("2006-01-02"))
	}
}
