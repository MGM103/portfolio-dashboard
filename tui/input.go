package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputFields struct {
	description string
	focusIndex  int
	inputs      []textinput.Model
}

func NewInputFields(numFields int, placeholders []string) inputFields {
	fields := inputFields{inputs: make([]textinput.Model, numFields)}

	var t textinput.Model
	for i := range fields.inputs {
		t = textinput.New()
		t.Width = 25

		if i < len(placeholders) {
			t.Placeholder = placeholders[i]
		}
		if i == 0 {
			t.Focus()
		}

		fields.inputs[i] = t
	}

	return fields
}

func (i inputFields) Init() tea.Cmd {
	return textinput.Blink
}

func (i inputFields) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()

		switch s {
		case "tab", "down":
			i.focusIndex++

		case "shift+tab", "up":
			i.focusIndex--
		}

		if i.focusIndex >= len(i.inputs) {
			i.focusIndex = 0
		}

		if i.focusIndex < 0 {
			i.focusIndex = len(i.inputs) - 1
		}

		cmds := make([]tea.Cmd, len(i.inputs))
		for index := range i.inputs {
			if index == i.focusIndex {
				cmds[index] = i.inputs[index].Focus()
			} else {
				i.inputs[index].Blur()
			}
		}
	}

	cmd := i.updateInputs(msg)

	return i, cmd
}

func (i *inputFields) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(i.inputs))

	for index := range i.inputs {
		i.inputs[index], cmds[index] = i.inputs[index].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (i inputFields) View() string {
	var b strings.Builder

	b.WriteString(i.description)
	b.WriteRune('\n')

	for index := range i.inputs {
		b.WriteString(i.inputs[index].View())
		if index < len(i.inputs) {
			b.WriteRune('\n')
		}
	}

	b.WriteString("\n\nPress <enter> to submit")

	return b.String()
}

func (i inputFields) GetValues() []string {
	values := make([]string, len(i.inputs))
	for index, input := range i.inputs {
		values[index] = input.Value()
	}

	return values
}
