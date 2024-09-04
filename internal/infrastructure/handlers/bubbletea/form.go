package bubbletea

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// FormModel is a model for form
type FormModel struct {
	inputs       []textinput.Model
	focusedInput int
	err          error

	// submit is a function which will be called after all inputs are filled
	// arguments are values of inputs in order they are stored in inputs slice
	submit func(values []string) error
}

// NewFormModel is a constructor for FormModel
func NewFormModel(inputs []textinput.Model, submit func(values []string) error) *FormModel {
	return &FormModel{
		inputs: inputs,
		submit: submit,
	}
}

// Init is an initialization function
func (m *FormModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update is an update function
func (m *FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focusedInput == len(m.inputs)-1 {
				values := make([]string, len(m.inputs))
				for i, input := range m.inputs {
					values[i] = input.Value()
				}
				if err := m.submit(values); err != nil {
					m.err = err
					return m, nil
				}
				for i := range m.inputs {
					m.inputs[i].SetValue("")
				}
				return m, tea.Quit
			}
			m.inputs[m.focusedInput].Blur()
			m.focusedInput++
			if m.focusedInput >= len(m.inputs) {
				m.focusedInput = 0
			}
			m.inputs[m.focusedInput].Focus()
		case tea.KeyTab, tea.KeyDown:
			m.inputs[m.focusedInput].Blur()
			m.focusedInput++
			if m.focusedInput >= len(m.inputs) {
				m.focusedInput = 0
			}
			m.inputs[m.focusedInput].Focus()
		case tea.KeyShiftTab, tea.KeyUp:
			m.inputs[m.focusedInput].Blur()
			m.focusedInput--
			if m.focusedInput < 0 {
				m.focusedInput = len(m.inputs) - 1
			}
			m.inputs[m.focusedInput].Focus()
		case tea.KeyCtrlC, tea.KeyEsc:
			for i := range m.inputs {
				m.inputs[i].SetValue("")
			}
			return m, tea.Quit
		default:
		}
	}
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

// View is a view function
func (m *FormModel) View() string {
	s := ""
	for i, input := range m.inputs {
		s += fmt.Sprintf("%s\n", input.View())
		if i == m.focusedInput {
			s += "\n"
		}
	}
	if m.err != nil {
		s += fmt.Sprintf("Error: %s\n", m.err.Error())
	}
	return s
}
