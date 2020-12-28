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
	"fmt"
	"github.com/getsolus/libeopkg/shared"
	"path/filepath"
	"strings"
)

// A Package is the Package section of the metadata file. It contains
// the main details that are important to users.
type Package struct {
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
	DistributionRelease int
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
func (p *Package) GetID() string {
	return filepath.Base(p.PackageURI)
}

// GetRelease is a helpful wrapper to return the package's current release
func (p *Package) GetRelease() int {
	return p.History[0].Release
}

// GetVersion is a helpful wrapper to return the package's current version
func (p *Package) GetVersion() string {
	return p.History[0].Version
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
func (p *Package) GetPathComponent() string {
	nom := strings.ToLower(p.Source.Name)
	letter := nom[0:1]
	if strings.HasPrefix(nom, "lib") && len(nom) > 3 {
		return  filepath.Join(nom[0:4], nom)
	}
    return filepath.Join(letter, nom)
}

// DeltaName returns the filename (without extension) of a delta package
func (p *Package) DeltaName(newRelease int) string {
	return fmt.Sprintf("%s-%d-%d-%d-%s",
		p.Name,
		p.GetRelease(),
		p.GetRelease(),
		p.DistributionRelease,
		p.Architecture)
}

// IsDeltaPossible will compare the two input packages and determine if it
// is possible for a delta to be considered. Note that we do not compare the
// distribution _name_ because Solus already had to do a rename once, and that
// broke delta updates. Let's not do that again. eopkg should in reality determine
// delta applicability based on repo origin + upgrade path, not names
func (p *Package) IsDeltaPossible(newPackage *Package) bool {
	return p.GetRelease() < newPackage.GetRelease() &&
		p.Name == newPackage.Name &&
		p.DistributionRelease == newPackage.DistributionRelease &&
		p.Architecture == newPackage.Architecture
}

// Clean removes unnecessary whitespace and adds missing language attributes
func (p *Package) Clean() {
	p.Summary.Clean()
	p.Description.Clean()
}
