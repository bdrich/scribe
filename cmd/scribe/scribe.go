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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bdrich/scribe/pkg/fileSystem"
	"github.com/bdrich/scribe/pkg/templates"
)

func main() {

	pwd, _ := os.Getwd()
	const usage = "USAGE:\tscribe \n\t\t-in=<input_path>\n\t\t-out=<output_path>\n\t\t-templates=<dir_path>\n"

	// flags
	inputFileName := flag.String("in", "", "Input Values File -- input file which populates templates")
	outputPathName := flag.String("out", "", "Output Path -- path with populated templates")
	templatesPath := flag.String("templates", "", "Template File(s) -- path to recursively populate")
	flag.Parse()

	if *outputPathName == "" || *inputFileName == "" || *templatesPath == "" {
		fmt.Println("Error: Input Values File, Output, Templates paths cannot be empty.")
		fmt.Printf("----------------------------------------\n%s\n", usage)
		os.Exit(1)
	}

	outputPath := pwd + "/" + *outputPathName
	inputPath := pwd + "/" + *inputFileName
	fmt.Printf("Input Values File: %s\n", inputPath)

	// Ensure input file exists
	if !fileSystem.PathExists(inputPath) {
		fmt.Printf("Error: Input Values File '%s' does not exist.\n", inputPath)
		return
	}

	// If env directory doesn't exist, create it
	if !fileSystem.PathExists(outputPath) {
		fileSystem.CreateDir(outputPath)
	}

	templates.WalkExecuteTemplates(*templatesPath, inputPath, outputPath)
}
