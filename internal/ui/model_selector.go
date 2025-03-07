package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/openai/openai-go"
)

type ModelSelector struct {
	choices  []string
	cursor   int
	selected string
	quit     bool
}

func NewModelSelector() *ModelSelector {
	return &ModelSelector{
		choices: []string{
			openai.ChatModelGPT4oMini,
			openai.ChatModelGPT4o,
			openai.ChatModelO3Mini,
		},
		cursor: 0,
		quit:   false,
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

	s += "\n\n"
	s += HelpStyle.Render("  Use ") + KeyStyle.Render("↑/↓") + HelpStyle.Render(" or ") + KeyStyle.Render("k/j") + HelpStyle.Render(" to navigate") + "\n"
	s += HelpStyle.Render("  Press ") + KeyStyle.Render("Enter") + HelpStyle.Render(" or ") + KeyStyle.Render("Space") + HelpStyle.Render(" to select") + "\n"
	s += HelpStyle.Render("  Press ") + KeyStyle.Render("q") + HelpStyle.Render(" or ") + KeyStyle.Render("Ctrl+C") + HelpStyle.Render(" to exit") + "\n"

	return s
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
