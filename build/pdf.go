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
package build

import (
	"os"
	"path"

	"github.com/zivlakmilos/author/data"
)

func buildPdf(project *data.Project) error {
	format := project.Format
	if format == "markdown" {
		format = "markdown+rebase_relative_paths"
	}

	args := []string{
		"-f", format,
		"-t", "pdf",
		"--template", path.Join(project.Pdf.Template, "template.tex"),
		"-s",
		"-o", path.Join(project.OutputFolder, project.Pdf.OutputFolder, project.Pdf.OutputFileName),
		"--listings",
		"-V lang=rs-SR",
		"--pdf-engine", "xelatex",
	}

	if len(project.Pdf.Args) > 0 {
		args = append(args, project.Pdf.Args...)
	}

	if project.TableOfContent {
		args = append(args, "--toc")
	}

	if len(project.Bibliography) > 0 {
		args = append(args,
			"--bibliography",
			project.Bibliography,
			"--citeproc",
		)
	}

	if project.Pdf.Biblatex {
		args = append(args, "--biblatex")
	}

	err := os.MkdirAll(path.Join(project.OutputFolder, project.Pdf.OutputFolder), os.ModePerm)
	if err != nil {
		return err
	}

	err = pandoc(project.Sources, args, timeout)
	if err != nil {
		return err
	}

	return nil
}
