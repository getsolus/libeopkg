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

package archive

import (
	"github.com/getsolus/libeopkg/shared"
	"path/filepath"
	"strings"
)

// A MetaPackage is the Package section of the metadata file. It contains
// the main details that are important to users.
type MetaPackage struct {
	Name string
	// Brief description, one line, of the package functionality
	Summary shared.LocalisedFields
	// A full fleshed description of the package
	Description         shared.LocalisedFields
	IsA                 string `xml:"IsA,omitempty"`    // Legacy
	PartOf              string `xml:"PartOf,omitempty"` // component
	License             []string
	RuntimeDependencies *[]shared.Dependency `xml:"RuntimeDependencies>Dependency,omitempty"`
	Conflicts           *[]string            `xml:"Conflicts>Package,omitempty"`
	Replaces            *[]string            `xml:"Replaces>Package,omitempty"`
	Provides            shared.Provides      `xml:"Provides,omitempty"`
	History             []shared.Update      `xml:"History>Update"`
	// Binary details
	BuildHost           string
	Distribution        string
	DistributionRelease string
	Architecture        string
	InstalledSize       int64
	PackageSize         int64
	PackageHash         string
	PackageURI          string
	PackageFormat       string // Version
	// Index needs this, so do we for source==release matching
	Source shared.Source
}

// GetID will return the package ID for ferryd
func (m *MetaPackage) GetID() string {
	return filepath.Base(m.PackageURI)
}

// GetRelease is a helpful wrapper to return the package's current release
func (m *MetaPackage) GetRelease() int {
	return m.History[0].Release
}

// GetVersion is a helpful wrapper to return the package's current version
func (m *MetaPackage) GetVersion() string {
	return m.History[0].Version
}

// Metadata contains all of the information a package can provide to a user
// prior to installation. This includes the name, version, release, and so
// forth.
//
// Every Package contains Metadata, and during eopkg indexing, a reduced
// version of the Metadata is emitted.
type Metadata struct {
	Source  shared.Source
	Package MetaPackage `xml:"Package"`
}

// GetPathComponent will get the source part of the string which is used
// in all subdirectories of the repository.
//
// For all packages with a source name of 4 or more characters, the path
// component will be split on this, i.e.:
//
//      libr/libreoffice
//
// For all other packages, the first letter of the source name is used, i.e.:
//
//      n/nano
//
func (m *MetaPackage) GetPathComponent() string {
	nom := strings.ToLower(m.Source.Name)
	letter := nom[0:1]
	var path string
	if strings.HasPrefix(nom, "lib") && len(nom) > 3 {
		path = filepath.Join(nom[0:4], nom)
	} else {
		path = filepath.Join(letter, nom)
	}
	return path
}
