package render

type Option func(*Engine)

// WithDefaultLayout sets the default layout for the engine
// if no layout is specified in the template this layout will be used.
// By default this is set to "app/layouts/application.html"
func WithDefaultLayout(layout string) Option {
	return func(e *Engine) {
		e.defaultLayout = layout
	}
}

// WithHelpers sets the helpers for the engine these helpers will be
// available in all templates rendered by this engine.
func WithHelpers(hps map[string]any) Option {
	return func(e *Engine) {
		for k, v := range hps {
			e.helpers[k] = v
		}
	}
}
