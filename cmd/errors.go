package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cowboy-bebug/kommit/internal/utils"
)

type CmdType int

const (
	InitCmd CmdType = iota
	RootCmd
	VersionCmd
)

var cmdErrorPrefix = map[CmdType]string{
	InitCmd:    "😰 Therapy session interrupted",
	RootCmd:    "😰 Commitment issues detected",
	VersionCmd: "No errors are returned from this command.",
}

func getErrorPrefix(cmd CmdType) string {
	return cmdErrorPrefix[cmd]
}

func HandleUnsupportedModelError(cmd CmdType, err error) {
	if errors.Is(err, utils.UnsupportedModelError{}) {
		fmt.Printf("%s: The therapist's qualification looks sus!\n", getErrorPrefix(cmd))
		fmt.Println("(Check your .kommitrc.yaml for supported models)")
		os.Exit(1)
	}
}
