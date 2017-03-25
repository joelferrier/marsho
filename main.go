package main

import (
	"github.com/mitchellh/cli"
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("marsho")

func main() {
	args := os.Args[1:]
	verbose := false
	// inject 'version' as the command if the --version flag is set
	for i, arg := range args {
		if arg == "-v" || arg == "-version" || arg == "--version" {
			verArgs := make([]string, len(args)+1)
			verArgs[0] = "version"
			copy(verArgs[1:], args)
			args = verArgs
			break
		}

		if arg == "--verbose" {
			verbose = true
			// strip --verbose from args
			args = args[:i+copy(args[i:], args[i+1:])]
		}
	}

	InitLogging(verbose)

	c := &cli.CLI{
		Args:       args,
		Commands:   Commands,
		HelpFunc:   helpFunc,
		HelpWriter: os.Stdout,
	}

	c.Run()
}
