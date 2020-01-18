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
	"encoding/xml"
	"os"
	"sort"
)

const (
	// DefaultMaintainerName is the catch-all name for Solus maintainers
	DefaultMaintainerName = "Solus Team"
	// DefaultMaintainerEmail is the catch-all email for Solus maintainers
	DefaultMaintainerEmail = "root@solus-project.com"
)

// A Component as seen through the eyes of XML
type Component struct {
	// ID of this component, i.e. "system.base"
	Name string
	// Translated short name
	LocalName []LocalisedField
	// Translated summary
	Summary []LocalisedField
	// Translated description
	Description []LocalisedField
	// Which group this component belongs to
	Group string
	// Maintainer for this component
	Maintainer struct {
		// Name of the component maintainer
		Name string
		// Contact e-mail address of component maintainer
		Email string // Contact e-mail address of component maintainer
	}
}

// Components is a simple helper wrapper for loading from components.xml files
type Components struct {
	// Components is a list of Components
	Components ComponentList `xml:"Components>Component"`
}

// ComponentList allows us to quickly sort our components by name
type ComponentList []Component

// Len returns the length of a ComponentList
func (l ComponentList) Len() int {
	return len(l)
}

// Less returns true if the name of the first component is a lower value
func (l ComponentList) Less(a, b int) bool {
	return l[a].Name < l[b].Name
}

// Swap exchanges two components for sorting
func (l ComponentList) Swap(a, b int) {
	l[a], l[b] = l[b], l[a]
}

// NewComponents will load the Components data from the XML file
func NewComponents(xmlfile string) (cs *Components, err error) {
	// Open the component file
	cFile, err := os.Open(xmlfile)
	if err != nil {
		return
	}
	defer cFile.Close()
	// Decode the file contents
	cs = &Components{}
	dec := xml.NewDecoder(cFile)
	if err = dec.Decode(cs); err != nil {
		return
	}
	// Sort components by name
	sort.Sort(cs.Components)
	// Ensure there are no empty Lang= fields
	for i := range cs.Components {
		comp := &cs.Components[i]
		FixMissingLocalLanguage(&comp.LocalName)
		FixMissingLocalLanguage(&comp.Summary)
		FixMissingLocalLanguage(&comp.Description)
	}
	return
}
