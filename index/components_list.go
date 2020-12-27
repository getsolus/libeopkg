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
	"encoding/xml"
	"os"
	"sort"
)

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
func (l ComponentList) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

// Swap exchanges two components for sorting
func (l ComponentList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
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
		comp.LocalName.FixMissingLocalLanguage()
		comp.Summary.FixMissingLocalLanguage()
		comp.Description.FixMissingLocalLanguage()
	}
	return
}
