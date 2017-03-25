package main

import "testing"

func TestInitLogging(*testing.T) {
	verbose := true
	InitLogging(verbose)
}
