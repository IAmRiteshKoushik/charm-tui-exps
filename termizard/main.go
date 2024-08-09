package main

import (
	"log"

    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)


type model struct {
    index       int
    questions   []string
    width       int
    height      int
    answerField textinput.Model
}

func New(questions []string) *model {
    answerField := textinput.New()
    return &model{
        questions: questions,
        answerField: answerField,
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m model) View() string {
    if m.width == 0 {
        return "loading ..."
    }
    return lipgloss.JoinVertical(lipgloss.Center, m.questions[m.index], m.answerField.View())
}

func main(){
    // Setting up basic questions
    questions := []string{
        "What is your name ?",
        "What is your favourite editor ?",
        "What is your favourite quote ?",
    } 
    m := New(questions)

    // Setting up debugging logs
    f, err := tea.LogToFile("debug.log", "DEBUG:")
    if err != nil {
        log.Fatal("err: %w", err)
    }
    defer f.Close()

    // Adding the questions slice to the program
    p  := tea.NewProgram(m, tea.WithAltScreen())
    // Running the application and check for errors
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}

