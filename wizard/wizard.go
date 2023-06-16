package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField   lipgloss.Style
}

func DefaultStyles() *Styles {
	stl := new(Styles)
	stl.BorderColor = lipgloss.Color("#E7E7E7")
	stl.InputField = lipgloss.NewStyle().BorderForeground(stl.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return stl
}

type model struct {
	index 		int
	questions 	[]Question
	width 		int
	height 		int
	answerField textinput.Model
	styles 		*Styles
}

type Question struct {
	question string
	answer   string
}

func NewQuestion(question string) Question {
	return Question{
		question: question,
	}
}

func New(questions []Question) *model {
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Placeholder = "type your answer here"
	answerField.Focus()
	return &model{
		questions: 		questions,
		answerField: 	answerField,
		styles: 		styles,
	}
}

func (mod model) Init() tea.Cmd {
	return nil
}

func (mod model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &mod.questions[mod.index]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		mod.width = msg.Width
		mod.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return mod, tea.Quit
		case "enter":
			current.answer = mod.answerField.Value()
			mod.answerField.SetValue("")
			log.Printf("question: %s, answer: %s",
				current.question,
				current.answer,
			)
			mod.Next()
			return mod, nil
		}
	}
	mod.answerField, cmd = mod.answerField.Update(msg)
	return mod, cmd
}

func (mod model) View() string {
	if mod.width == 0 {
		return "loading..."
	}
	return lipgloss.Place(
		mod.width,
		mod.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			mod.questions[mod.index].question,
			mod.styles.InputField.Render(mod.answerField.View()),
		),
	)
}

func (mod *model) Next() {
	if mod.index < len(mod.questions)-1 {
		mod.index++
	} else {
		mod.index = 0
	}
}

func main() {
	questions := []Question{
		NewQuestion("what is your name?"),
		NewQuestion("what is your favourite editor?"),
		NewQuestion("what is your favourite quote?"),
	}
	m := New(questions)

	fail, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %w", err)
	}
	defer fail.Close()
	play := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := play.Run(); err != nil {
		log.Fatal(err)
	}
}
