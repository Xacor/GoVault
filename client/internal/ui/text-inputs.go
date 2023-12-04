package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Xacor/go-vault/proto"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type editorModel struct {
	focusIndex int
	forms      map[string][]textinput.Model
	focusForm  string
	cursorMode cursor.Mode
}

func initialEditorModel(secret *proto.Secret) editorModel {
	m := editorModel{
		forms: map[string][]textinput.Model{
			typeCredentials: {},
			typeBankCard:    {},
			typeText:        {},
			typeBinary:      {},
		},
		focusForm: typeCredentials, // TODO: draw form by secret type
	}

	for i := range m.forms {
		switch i {
		case typeCredentials:
			nameInput := newNameInput()

			loginInput := textinput.New()
			loginInput.Cursor.Style = cursorStyle
			loginInput.CharLimit = 256
			loginInput.Placeholder = "Login"
			loginInput.PromptStyle = focusedStyle
			loginInput.TextStyle = focusedStyle

			passwordInput := textinput.New()
			passwordInput.Cursor.Style = cursorStyle
			passwordInput.CharLimit = 256
			passwordInput.Placeholder = "Password"
			passwordInput.PromptStyle = focusedStyle
			passwordInput.TextStyle = focusedStyle
			passwordInput.EchoMode = textinput.EchoPassword
			passwordInput.EchoCharacter = '•'

			m.forms[i] = []textinput.Model{nameInput, loginInput, passwordInput}
		case typeBankCard:
			nameInput := newNameInput()

			numberInput := textinput.New()
			numberInput.Cursor.Style = cursorStyle
			numberInput.CharLimit = 16
			numberInput.Placeholder = "Number"
			numberInput.PromptStyle = focusedStyle
			numberInput.TextStyle = focusedStyle

			expDateMonthInput := textinput.New()
			expDateMonthInput.Cursor.Style = cursorStyle
			expDateMonthInput.CharLimit = 2
			expDateMonthInput.Placeholder = "Month"
			expDateMonthInput.PromptStyle = focusedStyle
			expDateMonthInput.TextStyle = focusedStyle
			expDateMonthInput.Validate = func(s string) error {
				num, err := strconv.Atoi(s)
				if err != nil {
					return err
				}

				if num < 1 || num > 12 {
					return errors.New("expiration month not in range")
				}

				return nil
			}

			expDateYearInput := textinput.New()
			expDateYearInput.Cursor.Style = cursorStyle
			expDateYearInput.CharLimit = 2
			expDateYearInput.Placeholder = "Year"
			expDateYearInput.PromptStyle = focusedStyle
			expDateYearInput.TextStyle = focusedStyle
			expDateYearInput.Validate = func(s string) error {
				num, err := strconv.Atoi(s)
				if err != nil {
					return err
				}

				if num < 1 || num > 99 {
					return errors.New("expiration year not in range")
				}

				return nil
			}

			holderInput := textinput.New()
			holderInput.Cursor.Style = cursorStyle
			holderInput.CharLimit = 128
			holderInput.Placeholder = "Holder Name"
			holderInput.PromptStyle = focusedStyle
			holderInput.TextStyle = focusedStyle

			m.forms[i] = []textinput.Model{nameInput, expDateMonthInput, expDateYearInput, holderInput}

		case typeBinary:
			nameInput := newNameInput()

			binaryInput := textinput.New()
			binaryInput.Cursor.Style = cursorStyle
			binaryInput.Placeholder = "Binary data"
			binaryInput.PromptStyle = focusedStyle
			binaryInput.TextStyle = focusedStyle

			m.forms[i] = []textinput.Model{nameInput, binaryInput}
		case typeText:
			nameInput := newNameInput()

			textInput := textinput.New()
			textInput.Cursor.Style = cursorStyle
			textInput.Placeholder = "Text"
			textInput.PromptStyle = focusedStyle
			textInput.TextStyle = focusedStyle

			m.forms[i] = []textinput.Model{nameInput, textInput}
		}
	}

	return m
}

func newNameInput() textinput.Model {
	nameInput := textinput.New()
	nameInput.Cursor.Style = cursorStyle
	nameInput.CharLimit = 256
	nameInput.Placeholder = "Name"
	nameInput.PromptStyle = focusedStyle
	nameInput.TextStyle = focusedStyle

	return nameInput
}

func (m editorModel) Init() tea.Cmd {
	return nil
}

func (m editorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.forms[m.focusForm]))
			for i := range m.forms[m.focusForm] {
				cmds[i] = m.forms[m.focusForm][i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to current input
		case "tab":
			cmds = append(cmds, m.focus()...)

		// Set focus to next input
		case "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.forms[m.focusForm]) {
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.forms[m.focusForm]) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.forms[m.focusForm])
			}

			cmds = append(cmds, m.focus()...)

		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *editorModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.forms[m.focusForm]))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.forms[m.focusForm] {
		m.forms[m.focusForm][i], cmds[i] = m.forms[m.focusForm][i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m editorModel) View() string {
	var b strings.Builder

	for i := range m.forms[m.focusForm] {
		b.WriteString(m.forms[m.focusForm][i].View())
		if i < len(m.forms[m.focusForm])-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.forms[m.focusForm]) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func (m *editorModel) focus() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.forms[m.focusForm]))
	for i := 0; i <= len(m.forms[m.focusForm])-1; i++ {
		if i == m.focusIndex {
			// Set focused state
			cmds[i] = m.forms[m.focusForm][i].Focus()
			m.forms[m.focusForm][i].PromptStyle = focusedStyle
			m.forms[m.focusForm][i].TextStyle = focusedStyle
			continue
		}
		// Remove focused state
		m.forms[m.focusForm][i].Blur()
		m.forms[m.focusForm][i].PromptStyle = noStyle
		m.forms[m.focusForm][i].TextStyle = noStyle
	}
	return cmds
}
