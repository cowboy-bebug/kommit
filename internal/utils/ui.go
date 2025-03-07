package utils

import (
	"time"

	"github.com/briandowns/spinner"
)

func Spinner(suffix string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + suffix
	s.Color("fgGreen")
	return s
}
