package assets

import (
	"embed"

	"github.com/leapkit/leapkit/core/assets"
)

var (
	//go:embed *
	files embed.FS

	// Manager for the PathFor helper
	Manager = assets.NewManager(files)
)
