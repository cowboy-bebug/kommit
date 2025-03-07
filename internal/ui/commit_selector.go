package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type CommitOption string

const (
	CommitOptionProceed CommitOption = "proceed"
	CommitOptionEdit    CommitOption = "edit"
	CommitOptionRerun   CommitOption = "rerun"
	CommitOptionExit    CommitOption = "exit"
)

type CommitSelector struct {
	choices  []string
	labels   map[string]CommitOption
	cursor   int
	selected CommitOption
	quit     bool
}

func NewCommitSelector() *CommitSelector {
	choices := []string{
		"Yes, I'm ready to commit to this message " + KeyStyle.Render("(proceed)"),
		"Yes, but I need to edit it first " + KeyStyle.Render("(edit)"),
		"No, I need another therapy session for a better message " + KeyStyle.Render("(re-run)"),
		"No, I'm terminating this therapy session " + KeyStyle.Render("(quit)"),
	}

	labels := map[string]CommitOption{
		choices[0]: CommitOptionProceed,
		choices[1]: CommitOptionEdit,
		choices[2]: CommitOptionRerun,
		choices[3]: CommitOptionExit,
	}

	return &CommitSelector{
		choices: choices,
		labels:  labels,
		cursor:  0,
		quit:    false,
	}
}

func (m *CommitSelector) Init() tea.Cmd {
	return nil
}

func (m *CommitSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected = m.labels[m.choices[m.cursor]]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *CommitSelector) View() string {
	s := "\n" + TitleStyle.Render("🧐 Do you want to use this commit message?") + "\n"

	for i, choice := range m.choices {
		cursor := " "
		style := ItemStyle

		if i == m.cursor {
			cursor = ">"
			style = SelectedItemStyle
		}

		s += fmt.Sprintf("%s %s\n", cursor, style.Render(choice))
	}

	s += "\n\n"
	s += HelpStyle.Render("  Use ") + KeyStyle.Render("↑/↓") + HelpStyle.Render(" or ") + KeyStyle.Render("k/j") + HelpStyle.Render(" to navigate") + "\n"
	s += HelpStyle.Render("  Press ") + KeyStyle.Render("Enter") + HelpStyle.Render(" or ") + KeyStyle.Render("Space") + HelpStyle.Render(" to select") + "\n"
	s += HelpStyle.Render("  Press ") + KeyStyle.Render("q") + HelpStyle.Render(" or ") + KeyStyle.Render("Ctrl+C") + HelpStyle.Render(" to exit") + "\n"

	return s
}

func SelectCommit() (CommitOption, error) {
	model := NewCommitSelector()
	p := tea.NewProgram(model)

	_, err := p.Run()
	if err != nil {
		return "", err
	}

	if model.quit {
		return CommitOptionExit, QuitError{}
	}

	// Clear the help lines
	for range 3 {
		fmt.Print("\033[1A\033[2K")
	}

	return model.selected, nil
}
