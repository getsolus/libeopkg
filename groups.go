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

// A Group as seen through the eyes of XML
type Group struct {
	// ID of this group, i.e. "multimedia"
	Name string
	// Translated short name
	LocalName []LocalisedField
	// Display icon for this Group
	Icon string
}

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
func (l GroupList) Less(a, b int) bool {
	return l[a].Name < l[b].Name
}

// Swap exchanges groups while sorting
func (l GroupList) Swap(a, b int) {
	l[a], l[b] = l[b], l[a]
}

// NewGroups will load the Groups data from the XML file
func NewGroups(xmlfile string) (*Groups, error) {
	// Open the groups file
	fi, err := os.Open(xmlfile)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	// Decode the contents
	grp := &Groups{}
	dec := xml.NewDecoder(fi)
	if err = dec.Decode(grp); err != nil {
		return nil, err
	}
	// Sort the groups
	sort.Sort(grp.Groups)
	// Ensure there are no empty Lang= fields
	for i := range grp.Groups {
		group := &grp.Groups[i]
		FixMissingLocalLanguage(&group.LocalName)
	}
	return grp, nil
}
