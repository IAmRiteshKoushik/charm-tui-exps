package main

import (
    "log"
    "fmt"
    "os"

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
    styles      *Styles
    index       int
    questions   []Question
    width       int
    height      int
    // answerField textinput.Model -> deprecated
    done        bool
}

// Constructor for model struct
func New(questions []Question) *model {
    styles := DefaultStyles()

    // -> deprecated
    // answerField := textinput.New()
    // answerField.Placeholder = "Your answer here:"
    // answerField.Focus()

    return &model{
        questions:  questions,
        styles:     styles,
    }
}

// func (m model) Init() tea.Cmd {
//     return m.questions[m.index].input.Blink
// }

type Question struct {
    question    string
    answer      string
    input       Input
}

// Constructor for Question struct
func NewQuestion(question string) Question{
    return Question{
        question: question,
    }
}

func newShortQuestion(question string) Question {
    q := NewQuestion(question)
    field := NewShortAnswerField()
    q.input = field
    return q 
}

func newLongQuestion(question string) Question {
    q := NewQuestion(question)
    field := NewLongAnswerField()
    q.input = field
    return q 
}


func (m model) Init() tea.Cmd {
    return m.questions[m.index].input.Blink
    // return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
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
            if m.index == len(m.questions) - 1 {
                m.done = true
            }
            current.answer = current.input.Value()
            m.Next()
            return m, current.input.Blur
        }
    }
    current.input, cmd = current.input.Update(msg)
    return m, cmd
}

func (m model) View() string {
    current := m.questions[m.index]
    if m.done {
        var output string
        for _, q := range m.questions {
            output += fmt.Sprintf("%s: %s\n", q.question, q.answer)
        }
        return output
    }
    if m.width == 0 {
        return "loading ..."
    }

    // Adding location to the components
    return lipgloss.Place(
        m.width,
        m.height,
        lipgloss.Center,
        lipgloss.Center, 
        lipgloss.JoinVertical(
            lipgloss.Left, 
            current.question,
            m.styles.InputField.Render(current.input.View()),
            ),
        )
}

func (m *model) Next() {
    // Change index or rotate through if last element is reached
    if m.index < len(m.questions) - 1 {
        m.index++
    } else {
        m.index = 0
    }
}

func main(){
    // Setting up basic questions
    questions := []Question {
        newShortQuestion("What is your name ?"),
        newShortQuestion("What is your favourite editor ?"),
        newLongQuestion("What is your favourite quote ?"),
    } 
    init := New(questions)

    // Setting up debugging logs
    f, err := tea.LogToFile("debug.log", "DEBUG:")
    if err != nil {
        log.Fatal("err: %w", err)
        os.Exit(1)
    }
    defer f.Close()

    // Adding the questions slice to the program
    p  := tea.NewProgram(*init, tea.WithAltScreen())
    // Running the application and check for errors
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
}
