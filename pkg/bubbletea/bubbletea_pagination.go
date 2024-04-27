package bubbletea

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/golang/protobuf/proto"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	spin = false
	// Avoid fetching multiple times while still fetching
	fetchingBackward = false
	fetchingForward  = false
)

type pageModel struct {
	items     *[]proto.Message
	paginator paginator.Model
	spinner   spinner.Model
}

func newModel(initMsg []proto.Message) pageModel {
	p := paginator.New()
	p.PerPage = msgPerPage
	p.SetTotalPages(len(initMsg))

	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("56"))
	s.Spinner = spinner.Points

	return pageModel{
		paginator: p,
		spinner:   s,
		items:     &initMsg,
	}
}

func (m pageModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m pageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
		switch {
		case key.Matches(msg, m.paginator.KeyMap.NextPage):
			// If no more data, don't fetch again (won't show spinner)
			value, ok := batchLen[lastBatchIndex+1]
			if !ok || value != 0 {
				if (m.paginator.Page >= (lastBatchIndex+1)*pagePerBatch-prefetchThreshold) && !fetchingForward {
					fetchingForward = true
					cmd = fetchDataCmd(lastBatchIndex+1, forward)
				}
			}
		case key.Matches(msg, m.paginator.KeyMap.PrevPage):
			if m.paginator.Page-firstBatchIndex*pagePerBatch == 0 {
				return m, cmd
			}
			if (m.paginator.Page <= firstBatchIndex*pagePerBatch+prefetchThreshold) && (firstBatchIndex > 0) && !fetchingBackward {
				fetchingBackward = true
				cmd = fetchDataCmd(firstBatchIndex-1, backward)
			}
		}
	case newDataMsg:
		if msg.fetchDirection == forward {
			// If current page is not in the range of the last batch, don't update
			if m.paginator.Page/pagePerBatch >= lastBatchIndex {
				*m.items = append(*m.items, msg.newItems...)
				lastBatchIndex++
				// If the number of batches exceeds the limit, remove the first batch
				if lastBatchIndex-firstBatchIndex >= localBatchLimit {
					*m.items = (*m.items)[batchLen[firstBatchIndex]:]
					firstBatchIndex++
				}
			}
			fetchingForward = false
		} else {
			// If current page is not in the range of the first batch, don't update
			if m.paginator.Page/pagePerBatch <= firstBatchIndex {
				*m.items = append(msg.newItems, *m.items...)
				firstBatchIndex--
				// If the number of batches exceeds the limit, remove the last batch
				if lastBatchIndex-firstBatchIndex >= localBatchLimit {
					*m.items = (*m.items)[:len(*m.items)-batchLen[lastBatchIndex]]
					lastBatchIndex--
				}
			}
			fetchingBackward = false
		}
		m.paginator.SetTotalPages(countTotalPages())
		return m, nil
	}

	m.paginator, _ = m.paginator.Update(msg)

	return m, cmd
}

func (m pageModel) View() string {
	var b strings.Builder
	table, err := getTable(&m)
	if err != nil {
		return "Error rendering table"
	}
	b.WriteString(table)
	b.WriteString(fmt.Sprintf("  PAGE - %d   ", m.paginator.Page+1))
	if spin {
		b.WriteString(fmt.Sprintf("%s%s", m.spinner.View(), " Loading new pages..."))
	}
	b.WriteString("\n\n  h/l ←/→ page • q: quit\n")

	return b.String()
}

func Paginator(_listHeader []printer.Column, _callback DataCallback) {
	listHeader = _listHeader
	callback = _callback

	var msg []proto.Message
	for i := firstBatchIndex; i < lastBatchIndex+1; i++ {
		newMessages := getMessageList(i)
		msg = append(msg, newMessages...)
	}

	p := tea.NewProgram(newModel(msg))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
