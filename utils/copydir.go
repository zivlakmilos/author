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

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func CopyDir(f embed.FS, origin, target string) error {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		if err := os.MkdirAll(target, os.ModePerm); err != nil {
			err = fmt.Errorf("error creating directory: %v", err)
			return err
		}
	}

	files, err := f.ReadDir(origin)
	if err != nil {
		err = fmt.Errorf("error reading directory: %v", err)
		return err
	}

	for _, file := range files {
		sourceFileName := filepath.Join(origin, file.Name())
		destFileName := filepath.Join(target, file.Name())

		if file.IsDir() {
			if err := CopyDir(f, sourceFileName, destFileName); err != nil {
				err = fmt.Errorf("error copying subdirectory: %v", err)
				return err
			}
			continue
		}

		fileContent, err := f.ReadFile(sourceFileName)
		if err != nil {
			err = fmt.Errorf("error reading file: %v", err)
			return err
		}

		if err := os.WriteFile(destFileName, fileContent, 0644); err != nil { // nolint: gosec
			log.Printf("error os.WriteFile error: %v", err)
			err = fmt.Errorf("error writing file: %w", err)
			return err
		}
	}

	return nil
}

func IsDirExists(f embed.FS, origin string) bool {
	_, err := f.ReadDir(origin)
	if err != nil {
		ExitWithError(err)
		return false
	}

	return true
}
