package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cowboy-bebug/kommit/internal/models"
)

type ModelSelector struct {
	choices  []string
	cursor   int
	selected string
	quit     bool
}

func NewModelSelector() *ModelSelector {
	return &ModelSelector{
		choices: models.OpenAISupportedModels,
		cursor:  0,
		quit:    false,
	}
}

func (m *ModelSelector) Init() tea.Cmd {
	return nil
}

func (m *ModelSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.selected = m.choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *ModelSelector) View() string {
	s := TitleStyle.Render("🧐 Choose your therapist's qualifications:") + "\n"

	for i, choice := range m.choices {
		cursor := " "
		style := ItemStyle

		if i == m.cursor {
			cursor = ">"
			style = SelectedItemStyle
		}

		s += cursor + " " + style.Render(choice) + "\n"
	}

	return WrapWithKeyboardHelp(s, false)
}

func SelectModel() (string, error) {
	model := NewModelSelector()
	p := tea.NewProgram(model)

	_, err := p.Run()
	if err != nil {
		return "", err
	}

	if model.quit {
		return "", QuitError{}
	}

	// Clear the help lines
	for range 3 {
		fmt.Print("\033[1A\033[2K")
	}

	return model.selected, nil
}
