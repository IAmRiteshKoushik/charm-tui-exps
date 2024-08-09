package main

import (
	"log"

    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/lipgloss" 
    tea "github.com/charmbracelet/bubbletea"
)

type Styles struct {
    BorderColor lipgloss.Color
    InputField  lipgloss.Style
}

// Constructor for Styles struct
func DefaultStyles() *Styles {
    s := new(Styles)
    s.BorderColor = lipgloss.Color("36")
    s.InputField = lipgloss.NewStyle().
        BorderForeground(s.BorderColor).
        BorderStyle(lipgloss.NormalBorder()).
        Padding(1).
        Width(80)

    return s
}

type model struct {
    index       int
    questions   []Question
    width       int
    height      int
    answerField textinput.Model
    styles      *Styles
}

// Constructor for model struct
func New(questions []Question) *model {
    styles := DefaultStyles()
    answerField := textinput.New()
    answerField.Placeholder = "Your answer here:"
    answerField.Focus()

    return &model{
        questions:      questions,
        answerField:    answerField,
        styles:         styles,
    }
}

type Question struct {
    question    string
    answer      string
}

// Constructor for Question struct
func NewQuestion(question string) Question{
    return Question{
        question: question,
    }
}


func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    current := &m.questions[m.index]
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "enter":
            current.answer = m.answerField.Value()
            m.answerField.SetValue("")
            log.Printf("Question: %s, answer %s", current.question, current.answer)
            m.Next()
            return m, nil
        }
    }
    var cmd tea.Cmd
    m.answerField, cmd = m.answerField.Update(msg)
    return m, cmd
}

func (m model) View() string {
    // If the screen is stuck on loading by any chance then we can send out 
    // a loading message
    if m.width == 0 {
        return "loading ..."
    }

    return lipgloss.Place(
        m.width,
        m.height,
        lipgloss.Center,
        lipgloss.Center, 
        lipgloss.JoinVertical(
            lipgloss.Center, 
            m.questions[m.index].question, 
            m.styles.InputField.Render(m.answerField.View()),
            ),
        )
}

func (m *model) Next() {
    if m.index < len(m.questions) - 1 {
        m.index++
    } else {
        m.index = 0
    }
}

func main(){
    // Setting up basic questions
    questions := []Question {
        NewQuestion("What is your name ?"),
        NewQuestion("What is your favourite editor ?"),
        NewQuestion("What is your favourite quote ?"),
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
