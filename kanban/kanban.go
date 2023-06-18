package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 4

const (
	todo status = iota
	inProgress
	done
)

/* CUSTOM ITEM */

type Task struct {
	status			status
	title			string
	description		string
}  

// implement the list.Item interface
func (task Task) FilterValue() string {
	return task.title
}

func (task Task) Title() string {
	return task.title
}

func (task Task) Description() string {
	return task.description
}

/* MAIN MODEL */

type Model struct {
	loaded 		bool
	focused 	status
	lists 		[]list.Model
	err 		error
}

func New() *Model {
	return &Model{}
}

// TODO: call this on tea.WindowSizeMsg
func (mod *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height)
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
			mod.initLists(msg.Width, msg.Height)
			mod.loaded = true
		}
	}
	var cmd tea.Cmd
	mod.lists[mod.focused], cmd = mod.lists[mod.focused].Update(msg)
	return mod, cmd
}

func (mod Model) View() string {
	if mod.loaded {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			mod.lists[todo].View(),
			mod.lists[inProgress].View(),
			mod.lists[done].View(),
		)
	}else{
		return "loading..."
	}
}

func main() {
	mod := New()
	prog := tea.NewProgram(mod)
	if err := prog.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
