package hctx

import "context"

type Context interface {
	context.Context
	New() Context
	Has(key string) bool
	Set(key string, value interface{})
}
