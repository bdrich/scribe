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

package fileSystem

import (
	"os"
)

// PathExists returns true if given path exists, false otherwise
func PathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// CreateDir creates a directory at given path, permission mode option
func CreateDir(path string, _mode ...int) error {
	var mode = 0777
	if len(_mode) > 0 {
		mode = _mode[0]
	}
	fileMode := os.FileMode(uint32(mode))
	err := os.MkdirAll(path, fileMode)
	return err
}

// CreateFile creates a file at given path
func CreateFile(path string) (*os.File, error) {
	file, err := os.Create(path)
	return file, err
}

// DeleteFile deletes a file at given path
func DeleteFile(path string) error {
	err := os.Remove(path)
	return err
}
