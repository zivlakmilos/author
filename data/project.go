package data

import (
	"encoding/json"
	"os"
)

type ProjectHtml struct {
	OutputFolder string `json:"outputFolder,omitempty"`
	Template     string `json:"template,omitempty"`
}

type ProjectPdf struct {
	OutputFolder   string   `json:"outputFolder,omitempty"`
	Template       string   `json:"template,omitempty"`
	OutputFileName string   `json:"outputFileName,omitempty"`
	Args           []string `json:"args,omitempty"`
	Biblatex       bool     `json:"biblatex,omitempty"`
}

type Project struct {
	Name           string      `json:"name,omitempty"`
	Author         string      `json:"author,omitempty"`
	Version        string      `json:"version,omitempty"`
	Format         string      `json:"format,omitempty"`
	TableOfContent bool        `json:"toc,omitempty"`
	Bibliography   string      `json:"bibliography,omitempty"`
	Sources        []string    `json:"sources,omitempty"`
	Assets         []string    `json:"assets,omitempty"`
	OutputFolder   string      `json:"outputFolder,omitempty"`
	Targets        []string    `json:"targets,omitempty"`
	Html           ProjectHtml `json:"html,omitempty"`
	Pdf            ProjectPdf  `json:"pdf,omitempty"`
}

func LoadProject(filepath string) (*Project, error) {
	var project Project

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}
