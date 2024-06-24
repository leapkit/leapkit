package content

import (
	"testing"

	"github.com/leapkit/core/internal/helpers/helptest"
	"github.com/stretchr/testify/require"
)

func Test_Default(t *testing.T) {
	r := require.New(t)
	def := "notavailable"

	t.Run("NotPresent", func(t *testing.T) {
		hc := helptest.NewContext()
		res := WithDefault("nothere", def, hc)
		r.Equal(res, def)
	})

	t.Run("Present", func(t *testing.T) {
		hc := helptest.NewContext()
		hc.Set("there", "someval")
		res := WithDefault("there", def, hc)
		r.Equal(res, "someval")
	})

	t.Run("PresentNil", func(t *testing.T) {
		hc := helptest.NewContext()
		hc.Set("nil", nil)
		res := WithDefault("nil", def, hc)
		r.Equal(res, "notavailable")
	})

	t.Run("PresentPointer", func(t *testing.T) {
		hc := helptest.NewContext()
		a := "something"
		p := &a
		hc.Set("pointer", p)
		res := WithDefault("pointer", def, hc)
		r.Equal(res, p)
		r.NotEqual(res, a)
	})
}
