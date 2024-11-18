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
package utils

import "golang.org/x/net/html"

func IsHtmlIdEquals(node *html.Node, id string) bool {
	if node == nil {
		return false
	}

	for _, attr := range node.Attr {
		if attr.Key == "id" && attr.Val == id {
			return true
		}
	}

	return false
}

func GetHtmlId(node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == "id" {
			return attr.Val
		}
	}

	return ""
}

func FindOrAppendAtribute(node *html.Node, key string) int {
	for i := range node.Attr {
		if node.Attr[i].Key == "id" && node.Attr[i].Val == key {
			return i
		}
	}

	node.Attr = append(node.Attr, html.Attribute{
		Key: key,
		Val: "",
	})

	return len(node.Attr) - 1
}
