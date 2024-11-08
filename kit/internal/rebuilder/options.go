package rebuilder

// NoOpRunner is a runner that does nothing.
var NoOpRunner = func() {}

type Option func(*configuration)

// WatchExtension sets the extensions to watch for changes.
func WatchExtension(extensions ...string) Option {
	return func(c *configuration) {
		c.extensionsToWatch = extensions
	}
}

func WithRunner(fn func()) Option {
	return func(c *configuration) {
		c.runners = append(c.runners, fn)
	}
}
