package ui

import (
	"errors"
	"fmt"
	"os"
)

type QuitError struct{}

func (q QuitError) Error() string {
	return "quitting"
}

func HandleQuitError(err error) {
	if errors.Is(err, QuitError{}) {
		fmt.Println("\n🥹 See you next time!")
		os.Exit(0)
	}
}
