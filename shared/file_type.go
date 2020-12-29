//
// Copyright © 2017-2020 Solus Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package shared

// FileType indicates the kind of file
type FileType string

const (
	// FileConfig is a configuration file (i.e. something in /etc)
	FileConfig FileType = "config"
	// FileData is a data file (catch-all)
	FileData = "data"
	// FileDoc is a piece of documentation for the package or its contents
	FileDoc = "doc"
	// FileExecutable is an executable binary or script
	FileExecutable = "executable"
	// FileHeader is a header file for a compiled language (e.g. C/C++)
	FileHeader = "header"
	// FileInfo is something in /usr/share/info
	FileInfo = "info"
	// FileLibrary is a shared or static library
	FileLibrary = "library"
	// FileLocale is the locale data used to translate applications
	FileLocale = "localedata"
	// FileMan is a manpage used for the "man" help system
	FileMan = "man"
)
