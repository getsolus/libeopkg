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

package pspec

import (
	"encoding/xml"
	"os"

	"github.com/getsolus/libeopkg/shared"
)

// PSpec is the format of the new-er `pspec_x86_64.xml` files
type PSpec struct {
	XMLName  xml.Name `xml:"PISI"`
	Source   shared.Source
	Packages []Package `xml:"Package"`
	History  []Update  `xml:"History>Update"`
}

// Load reads the pspec from a file
func Load(path string) (p *PSpec, err error) {
	p = &PSpec{}
	xmlFile, err := os.Open(path)
	if err != nil {
		return
	}
	defer xmlFile.Close()

	dec := xml.NewDecoder(xmlFile)
	err = dec.Decode(p)
	return
}
