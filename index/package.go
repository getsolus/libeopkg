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

package index

import (
	"github.com/getsolus/libeopkg/shared"
)

// Package represents one of the packages available in a repo
type Package struct {
	Name                string
	Summary             shared.LocalisedField
	Description         shared.LocalisedField
	IsA                 string              `xml:",omitempty"`
	PartOf              string              `xml:",omitempty"`
	Licenses            []string            `xml:"License"`
	RuntimeDependencies []shared.Dependency `xml:",omitempty"`
	Replaces            *[]string           `xml:"Replaces>Package,omitempty"`
	Conflicts           *[]string           `xml:"Conflicts>Package,omitempty"`
	Provides            *shared.Provides    `xml:",omitempty"`
	History             []shared.Update     `xml:"History>Update"`
	BuildHost           string
	Distribution        string
	DistributionRelease int
	Architecture        string
	InstalledSize       int
	PackageSize         int
	PackageHash         string
	PackageURI          string
	DeltaPackages       *[]Delta `xml:"DeltaPackages>Delta,omitempty"`
	PackageFormat       string
	Source              shared.Source
}
