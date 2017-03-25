package main

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/cli"
	"sort"
	"strings"
)

func helpFunc(commands map[string]cli.CommandFactory) string {
	// Find the maximum command length for formatting
	maxCmdLen := 0
	for key, _ := range commands {
		if len(key) > maxCmdLen {
			maxCmdLen = len(key)
		}
	}

	helpText := fmt.Sprintf(`
Usage: marsho [--version] [--help] <command> [args]
%s
`, listCommands(commands, maxCmdLen))
	return strings.TrimSpace(helpText)
}

func listCommands(commands map[string]cli.CommandFactory, maxCmdLen int) string {
	var buf bytes.Buffer

	keys := make([]string, 0, len(commands))
	for key, _ := range commands {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	//TODO: error checking
	for _, key := range keys {
		commandFunc, _ := commands[key]
		command, _ := commandFunc()
		key = fmt.Sprintf("%s%s", key, strings.Repeat(" ", maxCmdLen-len(key)))
		buf.WriteString(fmt.Sprintf("    %s    %s\n", key, command.Synopsis()))
	}
	return buf.String()
}
