package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

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
	list list.Model
	err error
}

func New() *Model {
	return &Model{}
}

// TODO: call this on tea.WindowSizeMsg
func (mod *Model) initList(width, height int) {
	mod.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	mod.list.Title = "To Do"
	mod.list.SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, and rice"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly clothes :)"},
	})
}

func (mod Model) Init() tea.Cmd {
	return nil
}

func (mod Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		mod.initList(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	mod.list, cmd = mod.list.Update(msg)
	return mod, cmd
}

func (mod Model) View() string {
	return mod.list.View()
}

func main() {
	mod := New()
	prog := tea.NewProgram(mod)
	if err := prog.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
