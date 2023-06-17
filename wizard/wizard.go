package main

import (
	"fmt"
	"log"
	"os"

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
	styles 		*Styles
	done 		bool
}

type Question struct {
	question string
	answer   string
	input	 Input
}

func NewQuestion(question string) Question {
	return Question{
		question: question,
	}
}

func newShortQuestion(question string) Question {
	question2 := NewQuestion(question)
	field := NewShortAnswerField()
	question2.input = field
	return question2
}

func newLongQuestion(question string) Question {
	question2 := NewQuestion(question)
	field := NewLongAnswerField()
	question2.input = field
	return question2
}

func New(questions []Question) *model {
	styles := DefaultStyles()
	return &model{
		questions: 		questions,
		styles: 		styles,
	}
}

func (main model) Init() tea.Cmd {
	return nil
}

func (main model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &main.questions[main.index]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		main.width = msg.Width
		main.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return main, tea.Quit
		case "enter":
			if main.index == len(main.questions)-1 {
				main.done = true
			}
			current.answer = current.input.Value()
			main.Next()
			return main, current.input.Blur
		}
	}
	current.input, cmd = current.input.Update(msg)
	return main, cmd
}

func (main model) View() string {
	current := main.questions[main.index]
	if main.done {
		var output string
		for _, question := range main.questions {
			output += fmt.Sprintf("%s: %s\n", question.question, question.answer)
		}
		return output
	}
	if main.width == 0 {
		return "loading..."
	}
	return lipgloss.Place(
		main.width,
		main.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			main.questions[main.index].question,
			main.styles.InputField.Render(current.input.View()),
		),
	)
}

func (main *model) Next() {
	if main.index < len(main.questions)-1 {
		main.index++
	} else {
		main.index = 0
	}
}

func main() {
	questions := []Question{
		newShortQuestion("what is your name?"),
		newShortQuestion("what is your favourite editor?"),
		newLongQuestion("what is your favourite quote?"),
	}
	main := New(questions)

	fail, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	}
	defer fail.Close()
	play := tea.NewProgram(*main, tea.WithAltScreen())
	if _, err := play.Run(); err != nil {
		log.Fatal(err)
	}
}
