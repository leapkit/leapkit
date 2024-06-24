package plush

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/leapkit/core/internal/plush/ast"
	"github.com/leapkit/core/render/hctx"
)

// Making sure that the interface is implemented
var _ hctx.HelperContext = &HelperContext{}

// Helpers contains all of the default helpers for
// These will be available to all templates. You should add
// any custom global helpers to this list.
var Helpers = HelperMap{
	moot: &sync.Mutex{},
}

func init() {
	Helpers.Add("partial", PartialHelper)
}

// HelperContext is an optional last argument to helpers
// that provides the current context of the call, and access
// to an optional "block" of code that can be executed from
// within the helper.
type HelperContext struct {
	hctx.Context
	compiler *compiler
	block    *ast.BlockStatement
}

const helperContextKind = "HelperContext"

// Render a string with the current context
func (h HelperContext) Render(s string) (string, error) {
	return Render(s, h.Context)
}

// HasBlock returns true if a block is associated with the helper function
func (h HelperContext) HasBlock() bool {
	return h.block != nil
}

// Block executes the block of template associated with
// the helper, think the block inside of an "if" or "each"
// statement.
func (h HelperContext) Block() (string, error) {
	return h.BlockWith(h.Context)
}

// BlockWith executes the block of template associated with
// the helper, think the block inside of an "if" or "each"
// statement, but with it's own context.
func (h HelperContext) BlockWith(hc hctx.Context) (string, error) {
	ctx, ok := hc.(*Context)
	if !ok {
		return "", fmt.Errorf("expected *Context, got %T", hc)
	}

	octx := h.compiler.ctx
	defer func() { h.compiler.ctx = octx }()
	h.compiler.ctx = ctx

	if h.block == nil {
		return "", fmt.Errorf("no block defined")
	}

	i, err := h.compiler.evalBlockStatement(h.block)
	if err != nil {
		return "", err
	}

	bb := &bytes.Buffer{}
	h.compiler.write(bb, i)

	return bb.String(), nil
}
