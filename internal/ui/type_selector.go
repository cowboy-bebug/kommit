package ui

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type TypeItem struct {
	Name     string
	Selected bool
}

type TypeSelector struct {
	types  []TypeItem
	cursor int
	help   help.Model
	quit   bool
}

func NewTypeSelector() *TypeSelector {
	defaultTypes := []string{
		"build",
		"chore",
		"ci",
		"docs",
		"feat",
		"fix",
		"perf",
		"refactor",
		"revert",
		"style",
		"test",
	}

	types := make([]TypeItem, len(defaultTypes))
	for i, t := range defaultTypes {
		types[i] = TypeItem{Name: t, Selected: true}
	}

	return &TypeSelector{
		types:  types,
		cursor: 0,
		help:   help.New(),
		quit:   false,
	}
}

func (m *TypeSelector) Init() tea.Cmd {
	return nil
}

func (m *TypeSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.types)-1 {
				m.cursor++
			}
		case " ":
			m.types[m.cursor].Selected = !m.types[m.cursor].Selected
		case "enter", "return":
			return m, tea.Quit
		case "a":
			for i := range m.types {
				m.types[i].Selected = true
			}
		case "n":
			for i := range m.types {
				m.types[i].Selected = false
			}
		}
	}
	return m, nil
}

func (m *TypeSelector) View() string {
	s := TitleStyle.Render("Select commit types for your therapy plan:") + "\n"

	for i, item := range m.types {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		checked := "[ ]"
		if item.Selected {
			checked = "[x]"
		}

		checkboxStyle := UncheckedStyle
		nameStyle := ItemStyle

		if item.Selected {
			checkboxStyle = CheckedStyle
		}

		if i == m.cursor {
			nameStyle = SelectedItemStyle
		}

		s += cursor + " " + checkboxStyle.Render(checked) + " " + nameStyle.Render(item.Name) + "\n"
	}

	return WrapWithKeyboardHelp(s,
		WithStandardNavigation(),
		WithSelectionOptions(),
	)
}

func (m *TypeSelector) GetSelectedTypes() []string {
	var selected []string
	for _, item := range m.types {
		if item.Selected {
			selected = append(selected, item.Name)
		}
	}
	return selected
}

func SelectTypes() ([]string, error) {
	model := NewTypeSelector()
	p := tea.NewProgram(model)

	_, err := p.Run()
	if err != nil {
		return nil, err
	}

	if model.quit {
		return nil, QuitError{}
	}

	return model.GetSelectedTypes(), nil
}
