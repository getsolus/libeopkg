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
	"github.com/getsolus/libeopkg/shared"
	"os"
)

// A Distribution as seen through the eyes of XML
type Distribution struct {
	// Name of source to match source repos
	SourceName string
	// Translated description
	Description shared.LocalisedFields
	// Published version number for compatibility
	Version int
	// Type of repository (should always be main, really. Just descriptive)
	Type string
	// Name of the binary repository
	BinaryName string
	// Package names that are no longer supported
	Obsoletes []string `xml:"Obsoletes>Package"`

	// fast lookup of obsoletes
	obsmap map[string]bool
}

// NewDistribution will load the Distribution data from the XML file
func NewDistribution(xmlfile string) (dist *Distribution, err error) {
	// Open XML file
	fi, err := os.Open(xmlfile)
	if err != nil {
		return
	}
	defer fi.Close()
	// Decode contenst
	dist = &Distribution{
		obsmap: make(map[string]bool),
	}
	dec := xml.NewDecoder(fi)
	if err = dec.Decode(dist); err != nil {
		return
	}
	// Build map of obsoletes
	for _, p := range dist.Obsoletes {
		dist.obsmap[p] = true
	}
	return
}

// IsObsolete will allow quickly determination of whether the package name
// was marked obsolete and should be hidden from the index
func (d *Distribution) IsObsolete(id string) bool {
	return d.obsmap[id]
}
