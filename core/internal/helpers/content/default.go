package content

import "github.com/leapkit/leapkit/core/render/hctx"

// WithDefault returns the key if exists, otherwise it returns
// defaultValue passed.
func WithDefault(key string, defaultValue interface{}, help hctx.HelperContext) interface{} {
	value := help.Value(key)
	if value != nil {
		return value
	}

	return defaultValue
}
