package ui

import (
	"github.com/Xacor/go-vault/proto"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	typeCredentials = "Credentials"
	typeBankCard    = "BankCard"
	typeText        = "Text"
	typeBinary      = "Binary"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type listModel struct {
	list         list.Model
	inputs       textinput.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
}

func initialListModel(secrets []*proto.Secret) listModel {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	items := make([]list.Item, len(secrets))
	for i := 0; i < len(secrets); i++ {
		item := item{
			title: secrets[i].Name,
		}

		if secrets[i].Credentials != nil {
			item.description = typeCredentials
		} else if secrets[i].BankCard != nil {
			item.description = typeBankCard
		} else if secrets[i].Text != nil {
			item.description = typeText
		} else if secrets[i].Binary != nil {
			item.description = typeBinary
		}

		items[i] = item
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	secretesList := list.New(items, delegate, 0, 0)
	secretesList.Title = "Secrets"
	secretesList.Styles.Title = titleStyle
	secretesList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return listModel{
		list:         secretesList,
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}
}

func (m listModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.keys.insertItem):
			m.delegateKeys.remove.SetEnabled(true)

			//open editor here
			newItem := item{}

			insCmd := m.list.InsertItem(0, newItem)
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))
			return m, tea.Batch(insCmd, statusCmd)
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m listModel) View() string {
	return appStyle.Render(m.list.View())
}

// type list struct {
// 	items    []string
// 	cursor   int
// 	selected int
// 	footer   string
// }

// func InitialList(items []*proto.Secret) list {
// 	list := list{
// 		items:    make([]string, 0, len(items)),
// 		cursor:   0,
// 		selected: -1,
// 	}
// 	for i := range items {
// 		list.items = append(list.items, items[i].Name)
// 	}

// 	return list
// }

// func (l list) Init() tea.Cmd {
// 	return nil
// }

// func (l list) View() string {
// 	s := "What should we buy at the market?\n\n"

// 	// Iterate over our choices
// 	for i := range l.items {

// 		// Is the cursor pointing at this choice?
// 		cursor := " " // no cursor
// 		if l.cursor == i {
// 			cursor = ">" // cursor!
// 		}

// 		// Is this choice selected?
// 		checked := " " // not selected
// 		if l.selected == i {
// 			checked = "x" // selected!
// 		}

// 		// Render the row
// 		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, l.items[i])
// 	}

// 	// The footer
// 	s += fmt.Sprintf("\n%s\n", l.footer)

// 	// Send the UI for rendering
// 	return s
// }

// func (l list) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {

// 	case tea.KeyMsg:
// 		switch msg.String() {

// 		case "ctrl+c", "q":
// 			return l, tea.Quit

// 		case "ctrl+s":
// 			if l.selected == l.cursor {
// 				l.footer = "Saved"
// 			}

// 		case "up", "k":
// 			if l.cursor > 0 {
// 				l.cursor--
// 			}
// 		case "down", "j":
// 			if l.cursor < len(l.items)-1 {
// 				l.cursor++
// 			}

// 		case "enter", " ":
// 			if l.cursor == l.selected {
// 				// TODO: close editor
// 				l.selected = -1
// 				l.footer = "Editor closed"
// 			} else {
// 				// TODO: open editor to the right
// 				l.selected = l.cursor
// 				l.footer = "Editor opened"
// 			}
// 		}
// 	}
// 	return l, nil
// }
