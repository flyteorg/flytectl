package bubbletea

import (
	"fmt"
	"log"
	"strings"
	"sync"

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
	// Avoid fetching back and forward at the same time
	mutex sync.Mutex
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
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
		switch {
		case key.Matches(msg, m.paginator.KeyMap.NextPage):
			// If already tried to fetch data and there is no more data to fetch (real last batch), don't fetch again (won't show spinner)
			value, ok := batchLen[lastBatchIndex+1]
			if !ok || value != 0 {
				if (m.paginator.Page >= (lastBatchIndex+1)*pagePerBatch-prefetchThreshold) && m.paginator.Page/pagePerBatch == lastBatchIndex && !fetchingForward {
					// if fetchingBack {
					// 	mutex.Lock()
					// }
					fetchingForward = true
					cmd = fetchDataCmd(lastBatchIndex+1, forward)
				}
			}
		case key.Matches(msg, m.paginator.KeyMap.PrevPage):
			if m.paginator.Page-firstBatchIndex*pagePerBatch == 0 {
				return m, cmd
			}
			if (m.paginator.Page <= firstBatchIndex*pagePerBatch+prefetchThreshold) && m.paginator.Page/pagePerBatch == firstBatchIndex && (firstBatchIndex > 0) && !fetchingBackward {
				// if fetchingForward {
				// 	mutex.Lock()
				// }
				fetchingBackward = true
				cmd = fetchDataCmd(firstBatchIndex-1, backward)
			}
		}

	// case spinner.TickMsg:
	// 	m.spinner, cmd = m.spinner.Update(msg)
	// 	return m, cmd
	case dataMsg:
		if len(msg.newItems) != 0 {
			if msg.fetchDirection == forward {
				if m.paginator.Page/pagePerBatch >= lastBatchIndex {
					*m.items = append(*m.items, msg.newItems...)
					lastBatchIndex++
					if lastBatchIndex-firstBatchIndex >= localBatchLimit {
						fmt.Println("delete back")
						*m.items = (*m.items)[batchLen[firstBatchIndex]:]
						//fmt.Println(len(*m.items), firstBatchIndex, batchLen[firstBatchIndex], "after forward")
						//localPageIndex -= batchLen[firstBatchIndex] / msgPerPage
						firstBatchIndex++
					}
				}
				fetchingForward = false
			} else if msg.fetchDirection == backward {
				if m.paginator.Page/pagePerBatch <= firstBatchIndex {
					*m.items = append(msg.newItems, *m.items...)
					firstBatchIndex--
					if lastBatchIndex-firstBatchIndex >= localBatchLimit {
						fmt.Println("delete forward")
						*m.items = (*m.items)[:len(*m.items)-batchLen[lastBatchIndex]]
						// batchLen[lastBatchIndex] = 0
						lastBatchIndex--
					}
				}
				fetchingBackward = false
			}
			m.paginator.SetTotalPages(countTotalPages())
		}
		m.paginator.Page = _min(_max(m.paginator.Page, firstBatchIndex*pagePerBatch), (lastBatchIndex+1)*pagePerBatch-1)
		// localPageIndex = _max(m.paginator.Page-firstBatchIndex*pagePerBatch, 0)
		return m, nil
	}
	//fmt.Println(len(*m.items))

	// switch msg := msg.(type) {
	// case tea.KeyMsg:
	// 	switch {
	// 	case key.Matches(msg, m.paginator.KeyMap.PrevPage):
	// 		if localPageIndex == 0 {
	// 			return m, cmd
	// 		}
	// 	}
	// }

	// fmt.Println("before ", m.paginator.Page-firstBatchIndex*pagePerBatch, firstBatchIndex, lastBatchIndex)
	// fmt.Printf("before %v %d %d %d\n", batchLen, m.paginator.TotalPages, len(*m.items), m.paginator.Page)

	// switch msg := msg.(type) {

	m.paginator, _ = m.paginator.Update(msg)

	fmt.Println("after ", m.paginator.Page, firstBatchIndex, lastBatchIndex)
	fmt.Printf("after %v %d %d %d\n", batchLen, m.paginator.TotalPages, len(*m.items), m.paginator.Page)

	return m, cmd
}

type fetchDirection int

const (
	forward fetchDirection = iota
	backward
)

type dataMsg struct {
	newItems       []proto.Message
	batchIndex     int
	fetchDirection fetchDirection
}

func fetchDataCmd(batchIndex int, fetchDirection fetchDirection) tea.Cmd {
	// fmt.Println("fetching")
	spin = true
	return func() tea.Msg {

		msg := dataMsg{
			newItems:       getMessageList(batchIndex),
			batchIndex:     batchIndex,
			fetchDirection: fetchDirection,
		}
		// mutex.Unlock()
		spin = false
		return msg
	}
}

func (m pageModel) View() string {
	var b strings.Builder
	_, err := getTable(&m)
	if err != nil {
		return "Error rendering table"
	}
	// b.WriteString(table)
	// b.WriteString(fmt.Sprintf("  PAGE - %d   ", m.paginator.Page+1))
	// if spin {
	// 	b.WriteString(fmt.Sprintf("%s%s", m.spinner.View(), " Loading new pages..."))
	// }
	// b.WriteString("\n\n  h/l ←/→ page • q: quit\n")

	return b.String()
}

func Paginator(_listHeader []printer.Column, _callback DataCallback) {
	listHeader = _listHeader
	callback = _callback

	var msg []proto.Message
	for i := firstBatchIndex; i < lastBatchIndex+1; i++ {
		newMessages := getMessageList(i)
		msg = append(msg, newMessages...)
		// batchLen[i] = len(newMessages)
	}

	p := tea.NewProgram(newModel(msg))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
