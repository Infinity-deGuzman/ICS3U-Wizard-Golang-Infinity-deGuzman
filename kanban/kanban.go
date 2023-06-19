package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 3

const (
	todo status = iota
	inProgress
	done
)

/* MODEL MANAGEENT */
var models []tea.Model

const (
	model status = iota
	form
)

/* STYLES */
var (
	columnStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.HiddenBorder())
	focusedStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#241"))
)

/* CUSTOM ITEM */

type Task struct {
	status			status
	title			string
	description		string
}

func NewTask(status status, title, description string) Task {
	return Task{status: status, title: title, description: description}
}

func (ta *Task) Next() {
	if ta.status == done {
		ta.status = todo
	} else {
		ta.status++
	}
}

/*func (ta *Task) Prev() {
	if ta.status == todo {
		ta.status = done
	} else {
		ta.status--
	}
}*/

// implement the list.Item interface
func (ta Task) FilterValue() string {
	return ta.title
}

func (ta Task) Title() string {
	return ta.title
}

func (ta Task) Description() string {
	return ta.description
}

/* MAIN MODEL */

type Model struct {
	loaded 		bool
	focused 	status
	lists 		[]list.Model
	err 		error
	quitting 	bool
}

func New() *Model {
	return &Model{}
}

func (mod *Model) MoveToNext() tea.Msg {
	selectedItem := mod.lists[mod.focused].SelectedItem()
	selectedTask := selectedItem.(Task)
	mod.lists[selectedTask.status].RemoveItem(mod.lists[mod.focused].Index())
	selectedTask.Next()
	mod.lists[selectedTask.status].InsertItem(len(mod.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	return nil
}

// TODO: Go to next list
func (mod *Model) Next() {
	if mod.focused == done {
		mod.focused = todo
	} else {
		mod.focused++
	}
} 

// TODO: Go to prev list

func (mod *Model) Prev() {
	if mod.focused == todo {
		mod.focused = done
	} else {
		mod.focused--
	}
} 

// TODO: call this on tea.WindowSizeMsg
func (mod *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height/2)
	defaultList.SetShowHelp(false)
	mod.lists = []list.Model{defaultList, defaultList, defaultList}

	// Init To Do
	mod.lists[todo].Title = "To Do"
	mod.lists[todo].Title = "To Do"
	mod.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, and rice"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly clothes :)"},
	})

	// Init in progress
	mod.lists[inProgress].Title = "In Progress"
	mod.lists[inProgress].Title = "To Do"
	mod.lists[inProgress].SetItems([]list.Item{
		Task{status: todo, title: "write code", description: "don't worry, it's Go"},
	})

	// Init done
	mod.lists[done].Title = "Done"
	mod.lists[done].Title = "To Do"
	mod.lists[done].SetItems([]list.Item{
		Task{status: todo, title: "stay cool", description: "as a cucumber"},
	})
}

func (mod Model) Init() tea.Cmd {
	return nil
}

func (mod Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !mod.loaded {
			columnStyle.Width(msg.Width / divisor)
			focusedStyle.Width(msg.Width / divisor)
			columnStyle.Height(msg.Height - divisor)
			focusedStyle.Height(msg.Height - divisor)
			mod.initLists(msg.Width, msg.Height)
			mod.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			mod.quitting = true
			return mod, tea.Quit
		case "left", "h":
			mod.Prev()
		case "right", "l":
			mod.Next()
		case "enter":
			return mod, mod.MoveToNext
		case "n":
			models[model] = mod // save the current model
			models[form] = NewForm(mod.focused)
			return models[form].Update(nil)
		}
		case Task:
			task := msg
			return mod, mod.lists[task.status].
			InsertItem(
				len(mod.lists[task.status].Items()),
				task,
			)
	}
	var cmd tea.Cmd
	mod.lists[mod.focused], cmd = mod.lists[mod.focused].Update(msg)
	return mod, cmd
}

func (mod Model) View() string {
	if mod.quitting {
		return ""
	}
	if mod.loaded {
		todoView := mod.lists[todo].View()
		inProgView := mod.lists[inProgress].View()
		doneView := mod.lists[done].View()
		switch mod.focused {
		case inProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgView),
				columnStyle.Render(doneView),
			)
		case done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				columnStyle.Render(inProgView),
				focusedStyle.Render(doneView),
			)
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgView),
				columnStyle.Render(doneView),
			)
		}
	} else {
		return "loading..."
	}
}

/* FORM MODEL */
type Form struct {
	focused 	status
	title 		textinput.Model
	description textarea.Model
}

func NewForm(focused status) *Form {
	form := &Form{focused: focused}
	form.title = textinput.New()
	form.title.Focus()
	form.description = textarea.New()
	return form
}

func (mod Form) CreateTask() tea.Msg {
	// TODO: create a new task
	return NewTask(
		mod.focused,
		mod.title.Value(),
		mod.description.Value(),
	)
}

func (mod Form) Init() tea.Cmd {
	return nil
}

func (mod Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return mod, tea.Quit
		case "enter":
			if mod.title.Focused() {
				mod.title.Blur()
				mod.description.Focus()
				return mod, textarea.Blink
			} else {
				models[form] = mod
				return models[model], mod.CreateTask
			}
		}
	}
	if mod.title.Focused() {
		mod.title, cmd = mod.title.Update(msg)
		return mod, cmd
	} else {
		mod.description, cmd = mod.description.Update(msg)
		return mod, cmd
	}
}

func (mod Form) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		mod.title.View(),
		mod.description.View(),
	)
}

func main() {
	models = []tea.Model{New(), NewForm(todo)}
	mod := models[model]
	prog := tea.NewProgram(mod)
	if err := prog.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
