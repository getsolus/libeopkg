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

// Groups is a simple helper wrapper for loading from components.xml files
type Groups struct {
	Groups GroupList `xml:"Groups>Group"`
}

// GroupList allows us to quickly sort our groups by name
type GroupList []Group

// Len returns the size of the list for sorting
func (l GroupList) Len() int {
	return len(l)
}

// Less returns true if the name of Group A is less than Group B's
func (l GroupList) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

// Swap exchanges groups while sorting
func (l GroupList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// NewGroups will load the Groups data from the XML file
func NewGroups(xmlfile string) (gs *Groups, err error) {
	// Open the groups file
	f, err := os.Open(xmlfile)
	if err != nil {
		return
	}
	defer f.Close()
	// Decode the contents
	gs = &Groups{}
	dec := xml.NewDecoder(f)
	if err = dec.Decode(gs); err != nil {
		return
	}
	// Sort the groups
	sort.Sort(gs.Groups)
	// Ensure there are no empty Lang= fields
	for i := range gs.Groups {
		group := &gs.Groups[i]
		group.LocalName.FixMissingLocalLanguage()
	}
	return
}
