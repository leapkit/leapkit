package rebuilder

// Starts the application and listen for changes.
func Start(path string, options ...Option) error {
	config.path = path
	for _, option := range options {
		option(config)
	}

	changed := make(chan bool)
	go runManager(changed)
	go runWatcher(changed)

	for _, v := range config.runners {
		go v()
	}

	<-make(chan struct{})

	return nil
}
