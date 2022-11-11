package hbrowser

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ferama/pg/pkg/conf"
	"github.com/ferama/pg/pkg/history"
)

var (
	borderStyle = lipgloss.ThickBorder()

	style = lipgloss.NewStyle().
		BorderTop(true).
		BorderRight(true).
		BorderLeft(true).
		BorderBottom(true).
		BorderForeground(lipgloss.Color(conf.ColorFocus)).
		BorderStyle(borderStyle)

	itemStyle          = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle  = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	indexStyle         = lipgloss.NewStyle()
	indexStyleSelected = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
)

// https://github.com/charmbracelet/bubbletea/blob/master/examples/list-simple/main.go

type hBrowserStatesMsg struct {
}

type HBrowserSelectedMsg struct {
	Idx int
}

type listItem struct {
	Idx   int
	Value string
}

func (i listItem) FilterValue() string { return i.Value }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, itm list.Item) {
	i, ok := itm.(listItem)
	if !ok {
		return
	}

	if index == m.Index() {
		fmt.Fprint(w, indexStyleSelected.Render(fmt.Sprint(index+1)))
		fmt.Fprint(w, selectedItemStyle.Render(fmt.Sprint(i.Value)))
	} else {
		fmt.Fprint(w, indexStyle.Render(fmt.Sprint(index+1)))
		fmt.Fprint(w, itemStyle.Render(fmt.Sprint(i.Value)))
	}
}

type Model struct {
	err            error
	focused        bool
	terminalHeight int
	terminalWidth  int

	list list.Model
}

func New() *Model {

	m := &Model{
		err:     nil,
		focused: false,
	}
	m.setState()
	return m
}

func (m *Model) setState() tea.Msg {

	delegate := itemDelegate{}

	listModel := list.New(make([]list.Item, 0), delegate, 0, 0)
	listModel.DisableQuitKeybindings()
	listModel.Styles.Title.
		UnsetBackground().
		Underline(true).
		Foreground(lipgloss.Color(conf.ColorTitle))
	listModel.Styles.FilterPrompt.Foreground(lipgloss.Color(conf.ColorFocus))

	h := history.GetInstance()
	hitems := h.GetList()
	items := make([]list.Item, 0)
	for idx := len(hitems) - 1; idx >= 0; idx-- {
		i := hitems[idx]
		i = strings.ReplaceAll(i, "\n", " ")
		items = append(items, listItem{
			Idx:   idx,
			Value: i,
		})
	}

	listModel.SetDelegate(delegate)
	listModel.SetItems(items)
	listModel.Title = "Query History"

	m.list = listModel
	return hBrowserStatesMsg{}
}

func (m *Model) setDimensions() {
	style.Width(m.terminalWidth - 4)
	style.Height(m.terminalHeight - 2)

	m.list.SetSize(m.terminalWidth-4, m.terminalHeight-2)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Focus() tea.Cmd {
	m.focused = true
	return m.setState
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.setDimensions()

	case hBrowserStatesMsg:
		m.setDimensions()

	case tea.KeyMsg:
		if !m.focused {
			break
		}
		switch msg.Type {
		case tea.KeyBackspace:
			i := m.list.SelectedItem()
			idx := i.(listItem).Idx
			h := history.GetInstance()
			h.DeleteAtIdx(idx)

			cmd = m.setState
			cmds = append(cmds, cmd)
		case tea.KeyEnter:
			i := m.list.SelectedItem()
			idx := i.(listItem).Idx

			cmd := func() tea.Msg {
				return HBrowserSelectedMsg{
					Idx: idx,
				}
			}
			cmds = append(cmds, cmd)
		}

	case error:
		m.err = msg
		return m, nil
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return style.Render(m.list.View())
}
