package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand"
	"os"
	"slices"
)

type model struct {
	selected    int // Int between 0 and 1 that stores the selected button. 0 for button on the left 1 for button on the right
	tried       []string
	currentCode string
}

type newCodeMsg string

var characters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func initialModel() model {
	return model{
		tried:       []string{},
		currentCode: randCode([]string{}),
	}
}

var selectedButtonStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("81")).Foreground(lipgloss.Color("0")).
	PaddingLeft(1).PaddingRight(1)

var unselectedButtonStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("245")).
	PaddingLeft(1).PaddingRight(1)

var appBorder = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("81")).
	Padding(1)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "a":
			if m.selected != 0 {
				m.selected--
			}
		case "right", "d":
			if m.selected != 1 {
				m.selected++
			}
		case "enter":
			if m.selected == 0 {
				m.tried = append(m.tried, m.currentCode)
				return m, newCode(m.tried)
			} else {
				return m, tea.Quit
			}
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	case newCodeMsg:
		m.currentCode = string(msg)
	}
	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("%v\nIs this your code?\n", m.currentCode)
	if m.selected == 0 {
		s += selectedButtonStyle.Render("No") + " " + unselectedButtonStyle.Render("Yes")
	} else {
		s += unselectedButtonStyle.Render("No") + " " + selectedButtonStyle.Render("Yes")
	}
	return appBorder.Render(s) + "\n"
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
func newCode(tried []string) tea.Cmd {
	return func() tea.Msg {
		return newCodeMsg(randCode(tried))
	}
}

func randCode(tried []string) string {
	codeSlice := make([]rune, 6)
	for i := range codeSlice {
		codeSlice[i] = characters[rand.Intn(len(characters))]
	}
	code := string(codeSlice)
	if slices.Contains(tried, string(codeSlice)) {
		code = randCode(tried)
	}
	return code
}
