package render

import (
	"html/template"
	"io"
	"io/fs"
	"sync"

	"github.com/gobuffalo/plush/v5"
)

// NewEngine builds the render engine based on the
// file system and the options passed to it.
func NewEngine(fs fs.FS, options ...Option) *Engine {
	e := &Engine{
		templates: fs,

		values:  make(map[string]any),
		helpers: make(template.FuncMap),

		defaultLayout: "app/layouts/application.html",
	}

	for _, option := range options {
		option(e)
	}

	return e
}

type Engine struct {
	templates     fs.FS
	defaultLayout string

	moot    sync.Mutex
	helpers template.FuncMap
	values  map[string]any
}

func (e *Engine) Set(key string, value any) {
	e.moot.Lock()
	defer e.moot.Unlock()

	e.values[key] = value
}

func (e *Engine) SetHelper(key string, value any) {
	e.moot.Lock()
	defer e.moot.Unlock()

	e.helpers[key] = value
}

func (e *Engine) HTML(w io.Writer) *Page {
	p := &Page{
		fs:     e.templates,
		writer: w,

		defaultLayout: e.defaultLayout,
	}

	ctx := plush.NewContext()
	for k, v := range e.values {
		ctx.Set(k, v)
	}

	for k, v := range e.helpers {
		ctx.Set(k, v)
	}

	ctx.Set("partialFeeder", func(name string) (string, error) {
		return p.open(name)
	})

	p.context = ctx

	return p
}

func (e *Engine) RenderHTML(template string, values map[string]any) (string, error) {
	ctx := plush.NewContext()
	for k, v := range e.values {
		ctx.Set(k, v)
	}

	for k, v := range e.helpers {
		ctx.Set(k, v)
	}

	for k, v := range values {
		ctx.Set(k, v)
	}

	t, err := e.templates.Open(template)
	if err != nil {
		return "", err
	}

	f, err := io.ReadAll(t)
	if err != nil {
		return "", err
	}

	return plush.Render(string(f), ctx)
}
