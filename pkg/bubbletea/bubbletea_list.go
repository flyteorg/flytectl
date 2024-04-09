package bubbletea

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const (
	listHeight   = 14
	defaultWidth = 20
)

var (
	// titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(0, 0, 0, 0)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := string(i)

	fn := itemStyle.Render

	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, _ := m.list.SelectedItem().(item)
			targetArgs = append(targetArgs, string(i))

			var isLeaf bool
			m.list, isLeaf = newList(string(i))
			if isLeaf {
				m.quitting = true
				rootCmd.SetArgs(targetArgs)
				return m, tea.Quit
			}

			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return quitTextStyle.Render("")
	}
	return "\n" + m.list.View()
}

func genList(items []list.Item) list.Model {
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	// l.Title = "Choose one of the commands"
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(false)
	// l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return l
}

func ifRunBubbleTea(_rootCmd cobra.Command) (*cobra.Command, bool, error) {
	cmd, flags, err := _rootCmd.Find(os.Args[1:])
	if err != nil {
		return cmd, false, err
	}

	err = _rootCmd.ParseFlags(flags)
	if err != nil {
		return nil, false, err
	}

	format, err := _rootCmd.Flags().GetString("format")
	if format != "bubbletea" || err != nil {
		return nil, false, err
	} else {
		return cmd, true, err
	}
}

func Bubbletea(_rootCmd *cobra.Command) error {

	rootCmd = _rootCmd
	targetArgs = os.Args[1:]

	currentCmd, run, err := ifRunBubbleTea(*rootCmd)
	if err != nil {
		return err
	} else if !run {
		return nil
	}

	items := generateSubCmdItems(currentCmd)

	l := genList(items)
	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return nil
}
