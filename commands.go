package main

import (
	"github.com/joelferrier/marsho/command"
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func init() {

	Commands = map[string]cli.CommandFactory{

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				BuildTime: BuildTime,
				Revision:  GitCommit,
				Version:   Version,
			}, nil
		},

		"fetch": func() (cli.Command, error) {
			return &command.FetchCommand{}, nil
		},

		"find": func() (cli.Command, error) {
			return &command.FindCommand{}, nil
		},

		"list": func() (cli.Command, error) {
			return &command.ListCommand{}, nil
		},
	}
}
