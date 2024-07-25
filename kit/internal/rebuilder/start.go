package rebuilder

import "fmt"

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
		go func() {
			// Wrapping the function inside a recover
			defer func() {
				var err error
				r := recover()
				if r == nil {
					return
				}

				switch t := r.(type) {
				case error:
					err = t
				case string:
					err = fmt.Errorf(t)
				default:
					err = fmt.Errorf("%+v", t)
				}

				fmt.Println(err)
			}()

			v()
		}()
	}

	<-make(chan struct{})

	return nil
}
