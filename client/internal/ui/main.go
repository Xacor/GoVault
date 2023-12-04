package ui

import (
	"github.com/Xacor/go-vault/proto"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState uint

const (
	listView sessionState = iota
	editorView
)

var modelStyle = lipgloss.NewStyle().

	// Align(lipgloss.Center, lipgloss.Center).
	BorderStyle(lipgloss.NormalBorder())

type mainModel struct {
	title  string
	cols   []column
	state  sessionState
	list   listModel
	editor editorModel
}

func NewMainModel(secrets []*proto.Secret) mainModel {
	m := mainModel{
		state:  listView,
		list:   InitialListModel(secrets),
		editor: initialEditorModel(secrets[0]),
	}

	return m
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.list.Init(), m.editor.Init())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// TODO: promt to save changes
			return m, tea.Quit
		case "tab":
			if m.state == listView {
				m.state = editorView
			} else {
				m.state = listView
			}
		case "n":
			if m.state == listView {
				list, cmd := m.list.Update(msg)
				m.list = list.(listModel)
				m.state = editorView
				cmds = append(cmds, cmd)
			} else {
				cmds = append(cmds, m.editor.focus()...)
			}
		}
		switch m.state {
		// Update whichever model is focused
		case listView:
			m.list.list, cmd = m.list.list.Update(msg)
			cmds = append(cmds, cmd)
		case editorView:
			editor, cmd := m.editor.Update(msg)
			m.editor = editor.(editorModel)

			cmds = append(cmds, cmd)
		}

	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	// return m.list.View()
	return lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.editor.View())
}
