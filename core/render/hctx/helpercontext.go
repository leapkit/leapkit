package hctx

type HelperContext interface {
	Context
	Block() (string, error)
	BlockWith(Context) (string, error)
	HasBlock() bool
	Render(s string) (string, error)
}
