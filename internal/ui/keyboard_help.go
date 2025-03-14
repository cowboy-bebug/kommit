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
	Navigate = fmt.Sprintf("  %s %s %s %s %s\n", Use, NavigateKey1, Or, NavigateKey2, ToNavigate)
	Proceed  = fmt.Sprintf("  %s %s %s\n", Press, ProceedKey, ToProceed)
	Exit     = fmt.Sprintf("  %s %s %s %s %s\n", Press, ExitKey1, Or, ExitKey2, ToExit)
	Toggle   = fmt.Sprintf("  %s %s %s\n", Press, ToggleKey, ToToggle)
	Select   = fmt.Sprintf("  %s %s %s\n", Press, SelectKey, ToSelect)
	Deselect = fmt.Sprintf("  %s %s %s\n", Press, DeselectKey, ToDeselect)
)

type HelpOption func(*helpConfig)
type helpConfig struct {
	showNavigate bool
	showProceed  bool
	showExit     bool
	showToggle   bool
	showSelect   bool
	showDeselect bool
}

func WithNavigation() HelpOption {
	return func(c *helpConfig) {
		c.showNavigate = true
	}
}

func WithProceed() HelpOption {
	return func(c *helpConfig) {
		c.showProceed = true
	}
}

func WithExit() HelpOption {
	return func(c *helpConfig) {
		c.showExit = true
	}
}

func WithToggle() HelpOption {
	return func(c *helpConfig) {
		c.showToggle = true
	}
}

func WithSelect() HelpOption {
	return func(c *helpConfig) {
		c.showSelect = true
	}
}

func WithDeselect() HelpOption {
	return func(c *helpConfig) {
		c.showDeselect = true
	}
}

func WithSelectionOptions() HelpOption {
	return func(c *helpConfig) {
		c.showToggle = true
		c.showSelect = true
		c.showDeselect = true
	}
}

func WithBasicNavigation() HelpOption {
	return func(c *helpConfig) {
		c.showNavigate = true
		c.showExit = true
	}
}

func WithStandardNavigation() HelpOption {
	return func(c *helpConfig) {
		c.showNavigate = true
		c.showProceed = true
		c.showExit = true
	}
}

func WrapWithKeyboardHelp(s string, options ...HelpOption) string {
	// Default configuration
	config := &helpConfig{
		showNavigate: true,
		showExit:     true,
	}

	// Apply all options
	for _, option := range options {
		option(config)
	}

	s += "\n\n"

	// NOTE: order of options is important
	if config.showNavigate {
		s += Navigate
	}

	if config.showProceed {
		s += Proceed
	}

	if config.showExit {
		s += Exit
	}

	if config.showToggle {
		s += Toggle
	}

	if config.showSelect {
		s += Select
	}

	if config.showDeselect {
		s += Deselect
	}

	return s
}
