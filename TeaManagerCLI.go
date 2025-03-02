package main

import (
	"fmt"
	"os"

	"TeaManagerCLI/task"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// State tracking for different modes
type state int

const (
	Normal state = iota
	EnteringTitle
	EnteringDueDate
)

type model struct {
	table   table.Model
	mode    state
	title   textinput.Model
	dueDate textinput.Model
	newTask table.Row
	lastTaskId int
	tasks []task.Task
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			err := task.SaveTasks(m.tasks)
			if err != nil {
				fmt.Println("Error saving tasks:", err)
			}
			return m, tea.Quit
		case "up":
			if m.mode == Normal {
				m.table.SetCursor(m.table.Cursor() - 1)
			}
		case "down":
			if m.mode == Normal {
				m.table.SetCursor(m.table.Cursor() + 1)
			}

		case "enter":
			if m.mode == Normal {
				index := m.table.Cursor()
				if index < 0 || index >= len(m.table.Rows()) {
					return m, nil
				}
				row := m.table.Rows()[index]
				status := row[2]
				if status == "✔️" {
					row[2] = "❌"
				} else {
					row[2] = "✔️"
				}
				m.table.SetRows(m.table.Rows())
			}
			if m.mode == EnteringTitle {
				lastTaskId := len(m.table.Rows()) + 1
				m.newTask = table.Row{fmt.Sprintf("%d", lastTaskId), m.title.Value(), "❌", ""}
				m.mode = EnteringDueDate
				m.dueDate.SetValue("") 
				m.dueDate.Focus()      
				return m, nil
			}
			if m.mode == EnteringDueDate {
				m.newTask[3] = m.dueDate.Value()
				rows := append(m.table.Rows(), m.newTask) // Insert new row
				m.table.SetRows(rows)
				m.tasks = append(m.tasks, task.Task{ // Update tasks slice
                    Id:      m.lastTaskId,
                    Title:   m.newTask[1],
                    Status:  false,
                    DueDate: m.newTask[3],
                })
				m.lastTaskId++
				m.mode = Normal // Return to normal mode
				return m, nil
			}

		case "delete":
			index := m.table.Cursor()
			if index < 0 || index >= len(m.table.Rows()) {
				return m, nil
			}
			rows := m.table.Rows()
			m.table.SetRows(append(rows[:index], rows[index+1:]...))
			m.tasks = append(m.tasks[:index], m.tasks[index+1:]...)
		case "ctrl+n":
			if m.mode == Normal {
				m.mode = EnteringTitle
				m.title.SetValue("") 
				m.title.Focus()      
				return m, textinput.Blink
			}
		}
	}
	if m.mode == EnteringTitle {
		var cmd tea.Cmd
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	}
	if m.mode == EnteringDueDate {
		var cmd tea.Cmd
		m.dueDate, cmd = m.dueDate.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.mode == EnteringTitle {
		return fmt.Sprintf("Enter Task Title: %s", m.title.View())
	}
	if m.mode == EnteringDueDate {
		return fmt.Sprintf("Enter Due Date (DD-MM-YYYY): %s", m.dueDate.View())
	}
	return m.table.View()
}

func main() {
	tasks, err := task.LoadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		os.Exit(1)
	}
	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Title", Width: 20},
		{Title: "Status", Width: 10},
		{Title: "Due Date", Width: 12},
	}
	var rows []table.Row
	for _, t := range tasks {
		status := "❌"
		if t.Status {
			status = "✔️"
		}
		rows = append(rows, table.Row{fmt.Sprintf("%d", t.Id), t.Title, status, t.DueDate})
	}
	totalTasks := len(rows)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	titleInput := textinput.New()
	titleInput.Placeholder = "Enter Task Title"
	titleInput.Focus()
	dueDateInput := textinput.New()
	dueDateInput.Placeholder = "Enter Due Date"
	p := tea.NewProgram(model{
		table:   t,
		mode:    Normal,
		title:   titleInput,
		dueDate: dueDateInput,
		lastTaskId: totalTasks,
		tasks: tasks,
	}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
