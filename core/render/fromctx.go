package render

import (
	"context"
)

// FromCtx returns the render engine from the context
// when its called it also adds any value in the valuer
// into the page, this is useful for the middlewares such as
// session to add values to the page.
func FromCtx(ctx context.Context) *Page {
	page := ctx.Value("renderer").(*Page)

	// Setting values from the valuer in the page
	vlr, ok := ctx.Value("valuer").(interface{ Values() map[string]any })
	if !ok {
		return page
	}

	for k, v := range vlr.Values() {
		page.Set(k, v)
	}

	return page
}

// FromCtx returns the render engine from the context,
// it assumes the render render engine has been set
// by the render middleware.
func EngineFromCtx(ctx context.Context) *Engine {
	return ctx.Value("renderEngine").(*Engine)
}
