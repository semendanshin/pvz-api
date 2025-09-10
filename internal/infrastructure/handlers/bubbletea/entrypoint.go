package bubbletea

import (
	tea "github.com/charmbracelet/bubbletea"
	"reflect"
)

// MyModel is a wrapper for tea.Model
type MyModel struct {
	Title string
	Model tea.Model
}

// EntryPointModel is a model for entry point
type EntryPointModel struct {
	choices          []MyModel
	cursor           int
	subProgramActive bool
}

// NewEntryPointModel is a constructor for EntryPointModel
func NewEntryPointModel(models []MyModel) *EntryPointModel {
	return &EntryPointModel{
		choices:          models,
		cursor:           0,
		subProgramActive: false,
	}
}

// Init is an initialization function
func (m *EntryPointModel) Init() tea.Cmd {
	return nil
}

func (m *EntryPointModel) cursorForward() {
	m.cursor++
	if m.cursor >= len(m.choices) {
		m.cursor = 0
	}
}

func (m *EntryPointModel) cursorBackward() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.choices) - 1
	}
}

func (m *EntryPointModel) quit() tea.Cmd {
	if !m.subProgramActive {
		return tea.Quit
	}
	m.subProgramActive = false
	return nil
}

func (m *EntryPointModel) handleKeyboard(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyDown:
		m.cursorForward()
	case tea.KeyUp:
		m.cursorBackward()
	case tea.KeyEnter:
		m.subProgramActive = true
	case tea.KeyCtrlC, tea.KeyEsc:
		return m.quit()
	default:
	}

	return nil
}

func (m *EntryPointModel) innerUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, m.handleKeyboard(msg)
	}

	return m, nil
}

// Update is an update function
func (m *EntryPointModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if !m.subProgramActive {
		_, cmd = m.innerUpdate(msg)
	} else {
		m.choices[m.cursor].Model, cmd = m.choices[m.cursor].Model.Update(msg)
		if cmd != nil {
			if reflect.ValueOf(cmd).Pointer() == reflect.ValueOf(tea.Quit).Pointer() {
				m.subProgramActive = false
				cmd = nil
			}
		}
	}

	return m, cmd
}

// View is a view function
func (m *EntryPointModel) View() string {
	if m.subProgramActive {
		return m.choices[m.cursor].Model.View()
	}

	s := ""

	s += "Выберите действие:\n"

	i := 0
	for _, model := range m.choices {
		if i == m.cursor {
			s += "> "
		} else {
			s += "  "
		}
		s += model.Title + "\n"
		i++
	}

	return s
}
