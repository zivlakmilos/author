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

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zivlakmilos/author/utils"
)

type model struct {
	cfg            Config
	txtProjectName textinput.Model
	lstTemplate    list.Model
	err            error
}

type quitWithErrorMsg struct {
	err error
}

func showTUI(cfg Config) {
	p := tea.NewProgram(initModel(cfg), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		utils.ExitWithError(err)
	}

	if m, ok := m.(model); ok {
		if m.cfg.ProjectName == "" || m.cfg.Template == "" {
			return
		}

		if m.err != nil {
			utils.ExitWithError(m.err)
		}
		utils.PrintSuccess(fmt.Sprintf("Project '%s' created", m.cfg.ProjectName))
	}
}

func initModel(cfg Config) model {
	txtProjectName := textinput.New()
	txtProjectName.Placeholder = "name"
	txtProjectName.Focus()

	items := getTemplatesList()

	lstTemplate := list.New(items, list.NewDefaultDelegate(), 0, 0)
	lstTemplate.Title = "Template"

	return model{
		cfg:            cfg,
		txtProjectName: txtProjectName,
		lstTemplate:    lstTemplate,
		err:            nil,
	}
}

func (m model) Init() tea.Cmd {
	if m.cfg.ProjectName == "" {
		return textinput.Blink
	}

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			if m.cfg.ProjectName == "" {
				m.cfg.ProjectName = m.txtProjectName.Value()
			} else if m.cfg.Template == "" {
				i, ok := m.lstTemplate.SelectedItem().(item)
				if !ok {
					return m, tea.Quit
				}

				m.cfg.Template = i.title
			}

			if m.cfg.ProjectName == "" || m.cfg.Template == "" {
				return m, nil
			}

			return m, m.createProjectCmd
		}
	case tea.WindowSizeMsg:
		m.lstTemplate.SetSize(msg.Width, msg.Height)
	case quitWithErrorMsg:
		m.err = msg.err
		return m, tea.Quit
	}

	var cmd tea.Cmd = nil

	if m.cfg.ProjectName == "" {
		m.txtProjectName, cmd = m.txtProjectName.Update(msg)
	} else if m.cfg.Template == "" {
		m.lstTemplate, cmd = m.lstTemplate.Update(msg)
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

	if m.cfg.Template == "" {
		return m.lstTemplate.View()
	}

	return ""
}

func (m model) createProjectCmd() tea.Msg {
	err := createProject(m.cfg)
	if err != nil {
		return quitWithErrorMsg{err: err}
	}

	return tea.Quit()
}
