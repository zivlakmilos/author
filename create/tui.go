/*
Copyright © 2024 Milos Zivlak

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package create

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cfg            Config
	txtProjectName textinput.Model
}

func showTUI(cfg Config) {
	p := tea.NewProgram(initModel(cfg))
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}

func initModel(cfg Config) model {
	txtProjectName := textinput.New()
	txtProjectName.Placeholder = "name"
	txtProjectName.Focus()

	return model{
		cfg:            cfg,
		txtProjectName: txtProjectName,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd = nil

	if m.cfg.ProjectName == "" {
		m.txtProjectName, cmd = m.txtProjectName.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.cfg.ProjectName == "" {
		return fmt.Sprintf(
			"Project Name\n\n%s\n\n%s",
			m.txtProjectName.View(),
			"(esc to quit)",
		)
	}

	return ""
}
