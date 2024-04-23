package bubbletea

// A simple program demonstrating the paginator component from the Bubbles
// component library.

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/lipgloss"
	"github.com/golang/protobuf/proto"

	tea "github.com/charmbracelet/bubbletea"
)

type pageModel struct {
	items     []proto.Message
	paginator paginator.Model
}

func newModel(initMsg []proto.Message) pageModel {
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = defaultMsgPerPage
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(initMsg))

	return pageModel{
		paginator: p,
		items:     initMsg,
	}
}

func (m pageModel) Init() tea.Cmd {
	return nil
}

func (m pageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.paginator, cmd = m.paginator.Update(msg)
	preFetchPage(&m)
	return m, cmd
}

func (m pageModel) View() string {
	var b strings.Builder
	start, end := m.paginator.GetSliceBounds(len(m.items))
	table, err := printTable(&m, start, end)
	if err != nil {
		return ""
	}
	b.WriteString(table)
	b.WriteString("  " + m.paginator.View())
	b.WriteString("\n\n  h/l ←/→ page • q: quit\n")
	return b.String()
}

func showPagination(initMsg []proto.Message) {
	p := tea.NewProgram(newModel(initMsg))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
