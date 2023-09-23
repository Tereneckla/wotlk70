package main

import (
	"github.com/Tereneckla/wotlk/cmd/wowsimcli/cmd"
	"github.com/Tereneckla/wotlk/sim"
)

func init() {
	sim.RegisterAll()
}

// Version information.
// This variable is set by the makefile in the release process.
var Version string

func main() {
	cmd.Execute(Version)
}
