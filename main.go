package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	files       []os.FileInfo
	currentPath string
	cursor      int
}

func (m *model) GoUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *model) GoDown() {
	if m.cursor < len(m.files)-1 {
		m.cursor++
	}
}

func initialModel() model {
	currentPath, _ := filepath.Abs(".")
	return model{
		files:       getFiles(currentPath),
		currentPath: currentPath,
		cursor:      0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "k", "up":
			m.GoUp()
		case "j", "down":
			m.GoDown()
		case "backspace":
			path, _ := filepath.Abs(filepath.Join(m.currentPath, ".."))

			m.files = getFiles(path)
			m.currentPath = path
			m.cursor = 0

		case "enter":
			curr := m.files[m.cursor]
			path := filepath.Join(m.currentPath, curr.Name())
			if curr.IsDir() {

				m.files = getFiles(path)
				m.currentPath = path
				m.cursor = 0

			} else {
				editor := os.Getenv("EDITOR")
				if len(editor) == 0 {
					// TODO: do something when the editor variable isn't set
				}
				return m, tea.ExecProcess(exec.Command(editor, path), nil)
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	content := fmt.Sprintf("In: %s (%d files)\n\n", m.currentPath, len(m.files))

	for i, file := range m.files {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		if file.IsDir() {
			content += fmt.Sprintf("%s 󰉋 %s\n", cursor, file.Name())
		} else {
			content += fmt.Sprintf("%s 󰈚 %s\n", cursor, file.Name())
		}
	}

	content += "\n(q to exit, backspace to exit current directory)"

	return content
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func getFiles(path string) (files []os.FileInfo) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Errror: Coudln't read directory")
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Println("Errror: Coudln't read directory")
		}
		files = append(files, info)
	}

	return
}
