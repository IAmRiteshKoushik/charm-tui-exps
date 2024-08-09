package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
    Blink()             tea.Msg
    Value()             string
    Blur()              tea.Msg
    Update(tea.Msg)     (Input, tea.Cmd)
    View()              string
}


/* SHORT ANSWER FIELD */
type ShortAnswerField struct {
    textinput textinput.Model
}


func NewShortAnswerField() *ShortAnswerField {
    ti := textinput.New()
    ti.Placeholder = "Your answer here:"
    ti.Focus()
    return &ShortAnswerField{ti}
}

func (sa *ShortAnswerField) Value() string {
    return sa.textinput.Value()
}

func (sa *ShortAnswerField) Blur() tea.Msg {
    return sa.textinput.Blur
}

func (sa *ShortAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
    var cmd tea.Cmd
    sa.textinput, cmd = sa.textinput.Update(msg)
    return sa, cmd
}

func (sa *ShortAnswerField) View() string {
    return sa.textinput.View()
}

func (sa *ShortAnswerField) Blink() tea.Msg {
    return textinput.Blink()
}

/* LONG ANSWER FIELD */
type LongAnswerField struct {
    textarea textarea.Model
}

func NewLongAnswerField() *LongAnswerField {
    ta := textarea.New()
    ta.Placeholder = "Your answer here:"
    ta.Focus()
    return &LongAnswerField{
        textarea: ta,
    }
}

func (la *LongAnswerField) Value() string {
    return la.textarea.Value()
}

func (la *LongAnswerField) Blur() tea.Msg {
    return la.textarea.Blur
}

func (la *LongAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
    var cmd tea.Cmd
    la.textarea, cmd = la.textarea.Update(msg)
    return la, cmd
}
func (la *LongAnswerField) Blink() tea.Msg {
    return textarea.Blink()
}

func (la *LongAnswerField) View() string {
    return la.textarea.View()
}
