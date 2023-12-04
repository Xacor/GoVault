package ui

import (
	"github.com/Xacor/go-vault/proto"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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
	toggleHelpMenu key.Binding
	insertItem     key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "add item"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type listModel struct {
	list list.Model
	keys *listKeyMap
}

func InitialListModel(secrets []*proto.Secret) listModel {
	var (
		listKeys = newListKeyMap()
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
	secretesList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	secretesList.Title = "Secrets"
	secretesList.Styles.Title = titleStyle
	secretesList.Paginator.PerPage = 5
	secretesList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.insertItem,
			listKeys.toggleHelpMenu,
		}
	}

	return listModel{
		list: secretesList,
		keys: listKeys,
	}
}

func (m listModel) Init() tea.Cmd {
	return nil
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
		case key.Matches(msg, m.keys.insertItem):

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
	return m.list.View()
}
