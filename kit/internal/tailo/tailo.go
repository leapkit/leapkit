package tailo

import (
	"os"

	"github.com/leapkit/leapkit/kit/internal/rebuilder"
	"github.com/paganotoni/tailo"
)

func BuildRunner(input, output, config string) func() {
	// if input or output files do not exist, return a
	// noOp function
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return rebuilder.NoOpRunner
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		return rebuilder.NoOpRunner
	}

	return tailo.WatcherFn(
		tailo.UseInputPath(input),
		tailo.UseOutputPath(output),
		tailo.UseConfigPath(config),
	)
}
