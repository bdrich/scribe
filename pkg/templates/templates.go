/*
	Copyright © 2019 IBM
	Copyright © 2019 Brian Richardson <brianthemathguy@gmail.com>
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	    http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package templates

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bdrich/scribe/pkg/fileSystem"

	"github.com/Masterminds/sprig"
	yaml "gopkg.in/yaml.v2"
)

// ExecutePathTemplates reads a YAML/JSON file from the valuesIn reader, uses it as values
// to populate the tplPath path.
func ExecutePathTemplates(valuesIn io.Reader, tplPath string) (string, error) {
	tpl := template.Must(template.New(tplPath).Funcs(sprig.TxtFuncMap()).Parse(tplPath))

	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, valuesIn)
	if err != nil {
		fmt.Printf("Failed to read standard input: %v\n", err)
		os.Exit(1)
	}

	var values map[string]interface{}
	err = yaml.Unmarshal(buf.Bytes(), &values)
	if err != nil {
		fmt.Printf("Failed to parse standard input: %v\n", err)
		os.Exit(1)
	}
	tplPathReader := strings.NewReader(tplPath)
	outputPath := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, tplPathReader)
	err = tpl.Execute(outputPath, values)
	if err != nil {
		fmt.Printf("Failed to parse standard input: %v\n", err)
		os.Exit(1)
	}
	outputPathStr := fmt.Sprintf("%s", outputPath)
	return outputPathStr, err
}

// ExecuteTemplates -- Reads a YAML/JSON file from the valuesIn, uses it as values
// to populate the tplFile templates and writes the executed templates to
// the out writer.
func ExecuteTemplates(valuesIn io.Reader, out io.Writer, tplFile string) error {
	templateFileReader, _ := ioutil.ReadFile(tplFile)
	templateText := string(templateFileReader)
	tpl := template.Must(template.New(tplFile).Funcs(sprig.TxtFuncMap()).Parse(templateText))

	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, valuesIn)
	if err != nil {
		fmt.Printf("Failed to read standard input: %v\n", err)
		os.Exit(1)
	}

	var values map[string]interface{}
	err = yaml.Unmarshal(buf.Bytes(), &values)
	if err != nil {
		fmt.Printf("Failed to parse standard input: %v\n", err)
		os.Exit(1)
	}

	err = tpl.Execute(out, values)
	if err != nil {
		fmt.Printf("Failed to parse standard input: %v\n", err)
		os.Exit(1)
	}
	return nil
}

// WalkExecuteTemplates is the main driver of the templating execution.
// Traverses the input templateRootPath to define directories and files
func WalkExecuteTemplates(templatesRootPath string, inputPath string, newRootPath string) {
	templatesRootPathStat, _ := os.Stat(templatesRootPath)
	templateFileList := []string{}
	switch {
	// if current path is a a directory
	case templatesRootPathStat.IsDir():
		_ = filepath.Walk(templatesRootPath, func(path string, f os.FileInfo, err error) error {
			templateFileList = append(templateFileList, path)
			return nil
		})
	case !templatesRootPathStat.IsDir():
		templateFileList = append(templateFileList, templatesRootPath)
	}

	// For each path in templates
	for _, templatePath := range templateFileList {
		inputPathReader, _ := os.OpenFile(inputPath, os.O_RDONLY, os.ModeAppend)
		templatePathStat, _ := os.Stat(templatePath)
		templatePathEnd := strings.TrimPrefix(templatePath, templatesRootPath)
		newPath := ""
		if templatePathEnd == "" {
			newPath = newRootPath
		} else {
			newPath = newRootPath + templatePathEnd
		}
		newPath, err := ExecutePathTemplates(inputPathReader, newPath)
		if err != nil {
			fmt.Printf("Failed to output paths: %v\n", err)
			os.Exit(1)
		}
		inputPathReader, _ = os.OpenFile(inputPath, os.O_RDONLY, os.ModeAppend)
		switch {
		case templatePathStat.IsDir():
			// create dir if it doesn't already exist in envPath
			if !fileSystem.PathExists(newPath) {
				fileSystem.CreateDir(newPath)
			}
		case !templatePathStat.IsDir():
			var outputFile *os.File
			if !fileSystem.PathExists(newPath) {
				outputFile, _ = fileSystem.CreateFile(newPath)
			} else {
				_ = fileSystem.DeleteFile(newPath)
				outputFile, _ = fileSystem.CreateFile(newPath)
			}
			ExecuteTemplates(inputPathReader, outputFile, templatePath)
			fmt.Printf("Output: %s\n\tFrom File: %s\n", newPath, templatePath)
		}
	}
}
