package main

import (
	"fmt"
	"net/http"
	"time"
    "os"

	tea "github.com/charmbracelet/bubbletea"
)

// We need to check whether this website is up or has the server crashed
const url = "https://charm.sh/"

type model struct {
    status  int
    err     error
}

func checkServer() tea.Msg {
     c := &http.Client{ Timeout: 10 * time.Second }
    res, err := c.Get(url)
    if err != nil {
        // Wrap the error we received and return it
        return errMsg{err}
    }
    // Return the HTTP status code
    return statusMsg(res.StatusCode)
}

type statusMsg int

type errMsg struct { err error }

// For messages that contain errors it's often handy to also implement 
// the error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

func (m model) Init() (tea.Cmd) {
    return checkServer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type){

    case statusMsg:
        // Save the status message that the server returned with to the model 
        // also, tell the bubbletea runtime to quit as there is nothing to do
        // A final view will still be rendered before quitting as the model 
        // updates the changes are reflected "reactively"
        m.status = int(msg)
        return m, tea.Quit
    case errMsg:
        // If there is an error -> update the model and tell the runtime that 
        // we would like to quit
        m.err = msg
        return m, tea.Quit
    case tea.KeyMsg:
        // Even though Ctrl + C exists, is it better to implement a quit key 
        // just in case something goes wrong and your users are not able to 
        // quit the program
        if msg.Type == tea.KeyCtrlC {
            return m, tea.Quit
        }
    }

    // If we happen to get any other messages then we do not do anything
    return m, nil
}

func (m model) View() string {
    // If there is an error, print it out and do not do anything else
    if m.err != nil {
        return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
    }

    // Tell the user we are doing something
    s := fmt.Sprintf("Checking %s ...", url)

    // When the server responds with a status, add it to the current line
    if m.status > 0 {
        s += fmt.Sprintf("\n%d %s!", m.status, http.StatusText(m.status))
    }

    // Send of whatever we came up with above for rendering
    return "\n" + s + "\n\n"
}

func main(){
    if _, err := tea.NewProgram(model{}).Run(); err != nil {
        fmt.Printf("Uh oh, there was an error: %v\n", err)
        os.Exit(1)
    }
}
