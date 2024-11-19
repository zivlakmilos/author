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
	"time"

	"github.com/zivlakmilos/author/data"
	"github.com/zivlakmilos/author/utils"
	"golang.org/x/net/html"
)

type processHtml struct {
	body    *html.Node
	section *html.Node
}

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

	err = postProcessHtml(project)
	if err != nil {
		return err
	}

	return nil
}

func copyHtmlAssets(dst string, html *data.ProjectHtml, assets []string) error {
	/*
	 * TODO: Optimise by copying only changed files?
	 */
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

func postProcessHtml(project *data.Project) error {
	filePath := path.Join(project.OutputFolder, project.Html.OutputFolder, "index.html")
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	node, err := html.Parse(f)
	f.Close()
	if err != nil {
		return err
	}

	p := processHtml{}
	p.postProcessHtmlNode(node)

	f, err = os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = html.Render(f, node)
	if err != nil {
		return err
	}

	return nil
}

func (p *processHtml) postProcessHtmlNode(node *html.Node) {
	if node.Type == html.ElementNode {
		id := utils.GetHtmlId(node)

		if id == "author-toc" {
			p.postProcessHtmlToc(node)
			return
		}

		if id == "author-body" {
			p.body = node
			p.section = nil
			p.postProcessHtmlBody(node)
			return
		}

		if id == "author-date" {
			p.postProcessHtmlDate(node, false)
		}

		if id == "author-copyright-year" {
			p.postProcessHtmlDate(node, true)
		}
	}

	for n := node.FirstChild; n != nil; n = n.NextSibling {
		p.postProcessHtmlNode(n)
	}
}

func (p *processHtml) postProcessHtmlToc(node *html.Node) {
	if node.Type == html.ElementNode {
		switch node.Data {
		case "ul":
			idx := utils.FindOrAppendAtribute(node, "class")
			node.Attr[idx].Val += "nav flex-column fixed-column"
		case "li":
			idx := utils.FindOrAppendAtribute(node, "class")
			node.Attr[idx].Val += "nav-item"
		case "a":
			idx := utils.FindOrAppendAtribute(node, "class")
			node.Attr[idx].Val += "nav-link"
		}
	}

	for n := node.FirstChild; n != nil; n = n.NextSibling {
		p.postProcessHtmlToc(n)
	}
}

func (p *processHtml) postProcessHtmlBody(node *html.Node) {
	p.body = node

	var nn *html.Node
	for n := node.FirstChild; n != nil; n = nn {
		nn = n.NextSibling
		p.postProcessHtmlSection(n)
	}
}

func (p *processHtml) postProcessHtmlDate(node *html.Node, onlyYrea bool) {
	data := node.FirstChild
	if data == nil {
		return
	}

	date, err := time.Parse("2006-01-02", data.Data)
	if err != nil {
		return
	}

	if onlyYrea {
		data.Data = date.Format("2006")
		return
	}

	data.Data = date.Format("02.01.2006.")
}

func (p *processHtml) postProcessHtmlSection(node *html.Node) {
	if node.Type == html.ElementNode {
		if node.Data == "h1" {
			p.postProcessHtmlCreateSection(node)
		}

		if p.section != nil && node.Data != "section" {
			p.postProcessHtmlSectionElement(node)
		}

		p.postProcessHtmlTag(node)
	}
}

func (p *processHtml) postProcessHtmlCreateSection(node *html.Node) {
	id := utils.GetHtmlId(node)
	node.Attr = nil

	if p.section != nil {
		n := &html.Node{
			Type: html.ElementNode,
			Data: "div",
			Attr: []html.Attribute{
				{
					Namespace: "",
					Key:       "class",
					Val:       "divider",
				},
			},
		}
		p.body.AppendChild(n)
	}

	n := &html.Node{
		Type: html.ElementNode,
		Data: "section",
		Attr: []html.Attribute{
			{
				Namespace: "",
				Key:       "id",
				Val:       id,
			},
		},
	}
	p.body.AppendChild(n)

	p.section = n
}

func (p *processHtml) postProcessHtmlSectionElement(node *html.Node) {
	p.body.RemoveChild(node)
	p.section.AppendChild(node)
}

func (p *processHtml) postProcessHtmlTag(node *html.Node) {
	if node.Type == html.ElementNode {
		switch node.Data {
		case "img":
			p.postProcessHtmlImg(node)
		case "p":
			p.postProcessHtmlParagraph(node)
		case "table":
			p.postProcessHtmlTable(node)
		}
	}

	for n := node.FirstChild; n != nil; n = n.NextSibling {
		p.postProcessHtmlTag(n)
	}
}

func (p *processHtml) postProcessHtmlImg(node *html.Node) {
	idx := utils.FindOrAppendAtribute(node, "style")
	node.Attr[idx].Val = "max-width: 100%;"
}

func (p *processHtml) postProcessHtmlParagraph(node *html.Node) {
	idx := utils.FindOrAppendAtribute(node, "style")
	node.Attr[idx].Val = "text-indent: 20px; text-align: justify;"
}

func (p *processHtml) postProcessHtmlTable(node *html.Node) {
	idx := utils.FindOrAppendAtribute(node, "class")
	node.Attr[idx].Val = "table table-bordered"
}
