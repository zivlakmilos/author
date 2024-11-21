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
package watch

import (
	"slices"
	"time"

	"github.com/zivlakmilos/author/build"
	"github.com/zivlakmilos/author/data"
	"github.com/zivlakmilos/author/utils"
)

type watcher struct {
	cfg     Config
	project *data.Project

	sources        map[string]time.Time
	projectModTime time.Time
}

func newWatcher(cfg Config) *watcher {
	w := &watcher{
		cfg:     cfg,
		sources: map[string]time.Time{},
	}

	return w
}

func (w *watcher) loadProject() error {
	project, err := data.LoadProject("project.json")
	if err != nil {
		return err
	}

	modTime, err := utils.GetFileModTime("project.json")
	if err != nil {
		return err
	}

	w.project = project
	w.projectModTime = modTime

	w.project.Targets = slices.DeleteFunc(w.project.Targets, func(el string) bool {
		if el == "html" && !w.cfg.Html {
			return true
		} else if el == "pdf" && !w.cfg.Pdf {
			return true
		}

		return false
	})

	w.sources = map[string]time.Time{}
	for _, src := range w.project.Sources {
		w.sources[src] = time.Time{}
	}

	return nil
}

func (w *watcher) run() {
	for {
		w.reloadProject()
		w.runBuild()
		time.Sleep(interval)
	}
}

func (w *watcher) reloadProject() {
	modTime, err := utils.GetFileModTime("project.json")
	if err != nil {
		return
	}

	if modTime.After(w.projectModTime) {
		err = w.loadProject()
		if err != nil {
			return
		}
	}
}

func (w *watcher) runBuild() {
	rebuild := false
	for key, val := range w.sources {
		modTime, _ := utils.GetFileModTime(key)
		if modTime.After(val) {
			w.sources[key] = modTime
			rebuild = true
		}
	}

	if rebuild {
		defer utils.PrintInfo("watch for changes...")

		utils.PrintInfo("build started")
		err := build.BuildProjectRun(w.project)
		if err != nil {
			utils.PrintError(err)
			return
		}

		utils.PrintSuccess("build finished")
	}
}
