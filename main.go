package main

import (
	"github.com/cowboy-bebug/kommitment/cmd"
)

var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date
	cmd.Execute()
}
