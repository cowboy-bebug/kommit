package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cowboy-bebug/kommit/internal/utils"
)

type Model struct {
	table table.Model
	quit  bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	s := TitleStyle.Render("💰 Kommit Financial Therapy Session 💰") +
		"\n\n" + m.table.View() + "\n"
	return WrapWithKeyboardHelp(s, WithBasicNavigation())
}

func NewTableModel(costs utils.Costs) Model {
	columns := []table.Column{
		{Title: "Repository", Width: 45},
		{Title: "Cost ($)", Width: 15},
	}

	rows := []table.Row{}
	var totalCost float64

	thisRepoIndex := 0
	repoName, _ := utils.GetRepoName()
	for repo, cost := range costs {
		rows = append(rows, table.Row{string(repo), fmt.Sprintf("%.5f", cost)})
		totalCost += float64(cost)

		if string(repo) == repoName {
			thisRepoIndex = len(rows) - 1
		} else {
			thisRepoIndex++
		}
	}
	rows = append(rows, table.Row{"TOTAL", fmt.Sprintf("%.5f", totalCost)})

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)
	t.SetCursor(thisRepoIndex)

	s := table.DefaultStyles()
	s.Header = TableHeaderStyle
	s.Selected = SelectedItemStyle.Padding(0, 0)
	t.SetStyles(s)

	return Model{table: t, quit: false}
}

func CostTableModel(costs utils.Costs) error {
	model := NewTableModel(costs)
	p := tea.NewProgram(model)

	_, err := p.Run()
	if err != nil {
		return err
	}

	if model.quit {
		return QuitError{}
	}

	return nil
}
