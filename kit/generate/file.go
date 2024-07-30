package generate

import (
	action "github.com/leapkit/leapkit/kit/generate/actions"
	"github.com/leapkit/leapkit/kit/generate/migrations"
)

type file interface {
	Generate() error
}

type Params struct {
	Kind string
	Name string
	Path []string
}

func New(p Params) file {
	switch p.Kind {
	case "action":
		return action.New(p.Path)
	case "migration":
		return migrations.New(p.Name)
	default:
		return nil
	}
}
