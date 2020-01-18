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

package libeopkg

import (
	"path/filepath"
	"strings"
)

// A Packager identifies the person responsible for maintaining the source
// package. In terms of ypkg builds, it will indicate the last person who
// made a change to the package, allowing a natural "blame" system to work
// much like git.
type Packager struct {
	Name  string
	Email string
}

// Source provides the information relating to the source package within
// each binary package.
// This source identifies one or more packages coming from the same origin,
// i.e they have the same *source name*.
type Source struct {
	Name     string
	Homepage string `xml:"Homepage,omitempty"`
	Packager Packager
}

// A Dependency has various attributes which help determine what needs to
// be installed when updating or installing the package.
type Dependency struct {
	Name string `xml:",chardata"`

	// Release based dependencies
	ReleaseFrom int `xml:"releaseFrom,attr,omitempty"`
	ReleaseTo   int `xml:"releaseTo,attr,omitempty"`
	Release     int `xml:"release,attr,omitempty"`

	// Version based dependencies
	VersionFrom string `xml:"versionFrom,attr,omitempty"`
	VersionTo   string `xml:"versionTo,attr,omitempty"`
	Version     string `xml:"version,attr,omitempty"`
}

// Action represents an action to take upon applying an update, such as restarting the system.
type Action struct {
	Value   string `xml:",chardata"`
	Package string `xml:"package,attr,omitempty"`
}

// An Update forms part of a package's history, describing the version, release,
// etc, for each release of the package.
type Update struct {
	Release int    `xml:"release,attr"`
	Type    string `xml:"type,attr,omitempty"`
	Date    string
	Version string
	Comment struct {
		Value string `xml:",cdata"`
	}
	Name struct {
		Value string `xml:",cdata"`
	}
	Email    string
	Requires *[]Action `xml:"Requires>Action,omitempty"`
}

// A COMAR script
type COMAR struct {
	Value  string `xml:",chardata"`
	Script string `xml:"script,attr,omitempty"`
}

// Provides defines special items that might be exported by a package
type Provides struct {
	COMAR       []COMAR  `xml:"COMAR,omitempty"`
	PkgConfig   []string `xml:"PkgConfig,omitempty"`
	PkgConfig32 []string `xml:"PkgConfig32,omitempty"`
}

// Delta describes a delta package that may be used for an update to save on bandwidth
// for the users.
//
// Delta upgrades are determined by placing the <DeltaPackages> section into the index, with
// each Delta listed with a releaseFrom. If the user is currently using one of the listed
// releaseFrom IDs in their installation, that delta package will be selected instead of the
// full package.
type Delta struct {
	ReleaseFrom int `xml:"releaseFrom,attr,omitempty"`
	PackageURI  string
	PackageSize int64
	PackageHash string
}

// A MetaPackage is the Package section of the metadata file. It contains
// the main details that are important to users.
type MetaPackage struct {
	Name string
	// Brief description, one line, of the package functionality
	Summary []LocalisedField
	// A full fleshed description of the package
	Description         []LocalisedField
	IsA                 string `xml:"IsA,omitempty"`    // Legacy
	PartOf              string `xml:"PartOf,omitempty"` // component
	License             []string
	RuntimeDependencies *[]Dependency `xml:"RuntimeDependencies>Dependency,omitempty"`
	Conflicts           *[]string     `xml:"Conflicts>Package,omitempty"`
	Replaces            *[]string     `xml:"Replaces>Package,omitempty"`
	Provides            *Provides     `xml:"Provides,omitempty"`
	History             []Update      `xml:"History>Update"`
	// Binary details
	BuildHost           string
	Distribution        string
	DistributionRelease string
	Architecture        string
	InstalledSize       int64
	PackageSize         int64
	PackageHash         string
	PackageURI          string
	DeltaPackages       *[]Delta `xml:"DeltaPackages>Delta,omitempty"`
	PackageFormat       string   // Version
	// Index needs this, so do we for source==release matching
	Source Source
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
	Source  Source
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
