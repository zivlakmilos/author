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
	"io/fs"
	"os"
	"path"

	"github.com/zivlakmilos/author/data"
)

func buildHtml(project *data.Project) error {
	args := []string{
		"-f", project.Format,
		"-t", "html",
		"--template", path.Join(project.Html.Template, "index.html"),
		"-s",
		"-o", path.Join(project.OutputFolder, project.Html.OutputFolder, "index.html"),
	}

	if len(project.Html.Args) > 0 {
		args = append(args, project.Html.Args...)
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

	if project.Biblatex {
		args = append(args, "--biblatex")
	}

	err := os.MkdirAll(path.Join(project.OutputFolder, project.Html.OutputFolder, "assets"), os.ModePerm)
	if err != nil {
		return err
	}

	err = copyHtmlAssets(path.Join(project.OutputFolder, project.Html.OutputFolder),
		&project.Html, project.Assets)
	if err != nil {
		return err
	}

	err = pandoc(project.Sources, args, timeout)
	if err != nil {
		return err
	}

	return nil
}

func copyHtmlAssets(dst string, html *data.ProjectHtml, assets []string) error {
	fs.WalkDir(os.DirFS(path.Join(html.Template, "public")), ".", func(pth string, dir fs.DirEntry, err error) error {
		err = os.RemoveAll(path.Join(dst, pth))
		if err != nil {
			return err
		}

		return nil
	})

	err := os.CopyFS(dst, os.DirFS(path.Join(html.Template, "public")))
	if err != nil {
		return err
	}

	dstAssets := path.Join(dst, "assets")
	for _, asset := range assets {
		err := os.CopyFS(dstAssets, os.DirFS(asset))
		if err != nil {
			return err
		}
	}

	return nil
}
