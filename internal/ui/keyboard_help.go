package ui

import "fmt"

// keybindings
var (
	NavigateKey1 = KeyStyle.Render("↑/↓")
	NavigateKey2 = KeyStyle.Render("k/j")
	ProceedKey   = KeyStyle.Render("Enter")
	ExitKey1     = KeyStyle.Render("Ctrl+c")
	ExitKey2     = KeyStyle.Render("q")
	ToggleKey    = KeyStyle.Render("Space")
	SelectKey    = KeyStyle.Render("a")
	DeselectKey  = KeyStyle.Render("n")
)

// words
var (
	Use        = HelpStyle.Render("Use")
	Press      = HelpStyle.Render("Press")
	Or         = HelpStyle.Render("or")
	ToNavigate = HelpStyle.Render("to navigate")
	ToProceed  = HelpStyle.Render("to proceed")
	ToExit     = HelpStyle.Render("to exit")
	ToToggle   = HelpStyle.Render("to toggle")
	ToSelect   = HelpStyle.Render("to select")
	ToDeselect = HelpStyle.Render("to deselect")
)

// help messages
var (
	Navigate = fmt.Sprintf("  %s %s %s %s %s\n", Use, NavigateKey1, Or, NavigateKey2, ToProceed)
	Proceed  = fmt.Sprintf("  %s %s %s\n", Press, ProceedKey, ToProceed)
	Exit     = fmt.Sprintf("  %s %s %s %s %s\n", Press, ExitKey1, Or, ExitKey2, ToExit)
	Toggle   = fmt.Sprintf("  %s %s %s\n", Press, ToggleKey, ToToggle)
	Select   = fmt.Sprintf("  %s %s %s\n", Press, SelectKey, ToSelect)
	Deselect = fmt.Sprintf("  %s %s %s\n", Press, DeselectKey, ToDeselect)
)

func WrapWithKeyboardHelp(s string, isSelect bool) string {
	s += "\n\n"
	s += Navigate + Proceed + Exit
	if isSelect {
		s += Toggle + Select + Deselect
	}
	return s
}
