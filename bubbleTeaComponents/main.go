package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"runtime"
)

type errMsg error
var selectedUser="noone";
var count=0;

var clear map[string]func() // create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) // Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") // Linux example, it's tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") // Windows example, it's tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

type model struct {
	viewport      viewport.Model
	messages      []string
	textarea      textarea.Model
	senderStyle   lipgloss.Style
	selectedIndex int   // Track the selected index in the list
	names         []string
	err           error
}

func main() {
	CallClear() // Clear the terminal screen

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initialModel() model {
	names := []string{"adam", "athul", "aromal", "mario"}
	vp := viewport.New(30, 5)
	vp.SetContent("\x1b[31mVault\x1b[0m\n" + strings.Join(names, "\n"))

	ta := textarea.New()
	// Remove the placeholder setting
	// ta.Placeholder = "Select user"
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:  ta,
		messages:      []string{},
		viewport:      vp,
		senderStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		selectedIndex: -1, // Initialize with no selection
		names:         names,
		err:           nil,
	}
}


func CallClear() {
	value, ok := clear[runtime.GOOS] // runtime.GOOS -> linux, windows, darwin, etc.
	if ok {                            // if we defined a clear func for that platform:
		value() // we execute it
	} else { // unsupported platform
		panic("Your platform is unsupported! I can't clear the terminal screen :(")
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			if count==1 {
				selectedUser=m.names[m.selectedIndex];
				
								fmt.Printf(m.textarea.Value())
								count=1;
							}
			if m.selectedIndex >= 0 && m.selectedIndex < len(m.names) &&count!=1 {
selectedUser=m.names[m.selectedIndex];

				fmt.Printf("Selected: %s\n",selectedUser)
				count=1;
			}
		
		case tea.KeyUp, tea.KeyDown:
			// Navigate through the list
			if msg.Type == tea.KeyUp {
				m.selectedIndex--
			} else {
				m.selectedIndex++
			}

			// Clamp the selected index within the bounds of the list
			if m.selectedIndex < 0 {
				m.selectedIndex = 0
			} else if m.selectedIndex >= len(m.names) {
				m.selectedIndex = len(m.names) - 1
			}

			// Update the viewport content to highlight the selected item
			content := fmt.Sprintf(strings.Join(highlightSelected(m.names, m.selectedIndex), "\n"))
			m.viewport.SetContent(content)
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

// highlightSelected highlights the selected item in the list
func highlightSelected(items []string, selectedIndex int) []string {
	result := make([]string, len(items))
	for i, item := range items {
		if i == selectedIndex {
			result[i] = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(item) // Highlight in red
		} else {
			result[i] = item
		}
	}
	return result
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}