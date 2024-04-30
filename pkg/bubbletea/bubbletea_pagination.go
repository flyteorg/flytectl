package bubbletea

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/flyteorg/flytectl/pkg/filters"
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
	p.Page = int(filter.Page) - 1
	p.SetTotalPages(getLocalLastPage())

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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.paginator.KeyMap.PrevPage):
			// If previous page will be out of the range of the first batch, don't update
			if m.paginator.Page == firstBatchIndex*pagePerBatch {
				return m, cmd
			}
		}
	}

	m.paginator, _ = m.paginator.Update(msg)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
		switch {
		case key.Matches(msg, m.paginator.KeyMap.NextPage):
			if (m.paginator.Page >= (lastBatchIndex+1)*pagePerBatch-prefetchThreshold) && !fetchingForward {
				// If no more data, don't fetch again (won't show spinner)
				value, ok := batchLen[lastBatchIndex+1]
				if !ok || value != 0 {
					fetchingForward = true
					cmd = fetchDataCmd(lastBatchIndex+1, forward)
				}
			}
		case key.Matches(msg, m.paginator.KeyMap.PrevPage):
			if (m.paginator.Page <= firstBatchIndex*pagePerBatch+prefetchThreshold) && (firstBatchIndex > 0) && !fetchingBackward {
				fetchingBackward = true
				cmd = fetchDataCmd(firstBatchIndex-1, backward)
			}
		}
	case newDataMsg:
		if msg.fetchDirection == forward {
			// Update if current page is in the range of the last batch
			// e.g. user left last batch while still fetching, then don't update
			if m.paginator.Page/pagePerBatch >= lastBatchIndex {
				*m.items = append(*m.items, msg.newItems...)
				lastBatchIndex++
				if lastBatchIndex-firstBatchIndex >= localBatchLimit {
					*m.items = (*m.items)[batchLen[firstBatchIndex]:]
					firstBatchIndex++
				}
			}
			fetchingForward = false
		} else {
			// Update if current page is in the range of the first batch
			// e.g. user left first batch while still fetching, then don't update
			if m.paginator.Page/pagePerBatch <= firstBatchIndex {
				*m.items = append(msg.newItems, *m.items...)
				firstBatchIndex--
				if lastBatchIndex-firstBatchIndex >= localBatchLimit {
					*m.items = (*m.items)[:len(*m.items)-batchLen[lastBatchIndex]]
					lastBatchIndex--
				}
			}
			fetchingBackward = false
		}
		m.paginator.SetTotalPages(getLocalLastPage())
		return m, nil
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

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

func Paginator(_listHeader []printer.Column, _callback DataCallback, _filter filters.Filters) error {
	listHeader = _listHeader
	callback = _callback
	filter = _filter
	filter.Page = int32(_max(int(filter.Page), 1))
	firstBatchIndex = (int(filter.Page) - 1) / pagePerBatch
	lastBatchIndex = firstBatchIndex

	var msg []proto.Message
	for i := firstBatchIndex; i < lastBatchIndex+1; i++ {
		newMessages := getMessageList(i)
		fmt.Println("newMessages", len(newMessages))
		if int(filter.Page)-(firstBatchIndex*pagePerBatch) > int(math.Ceil(float64(len(newMessages))/msgPerPage)) {
			return fmt.Errorf("the specified page has no data, please enter a valid page number")
		}
		msg = append(msg, newMessages...)
	}

	p := tea.NewProgram(newModel(msg))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}
