package lib

import (
	"os"

	"github.com/op/go-logging"
)

var Log *logging.Logger = logging.MustGetLogger("bittorrent")

func ConfigureLogging(verbose bool) {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	leveledBackend := logging.AddModuleLevel(backend)

	if verbose {
		leveledBackend.SetLevel(logging.DEBUG, "")
	} else {
		leveledBackend.SetLevel(logging.INFO, "")
	}

	logging.SetBackend(leveledBackend)
}
