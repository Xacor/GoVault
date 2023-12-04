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

type mainModel struct {
	title  string
	state  sessionState
	list   listModel
	editor editorModel
}

func NewMainModel(secrets []*proto.Secret) mainModel {
	m := mainModel{
		state:  listView,
		list:   initialListModel(secrets),
		editor: initialEditorModel(secrets[0]),
	}

	return m
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.list.Init(), m.editor.Init())
}

func (m mainModel) View() string {
	var s string
	s += lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.editor.View())
	return s
}

func (m *mainModel) Next() {
	if m.index == len(spinners)-1 {
		m.index = 0
	} else {
		m.index++
	}
}
