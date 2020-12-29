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
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/getsolus/libeopkg/shared"
	"io"
	"os"
	"path/filepath"
)

// Index is downloaded on a per-reprosity basis to provide information about the repository's:
// - Packages
// - Metadata
type Index struct {
	XMLName      xml.Name `xml:"PISI"`
	Distribution Distribution
	Packages     []Package   `xml:"Package"`
	Components   []Component `xml:"Component"`
	Groups       []Group     `xml:"Group"`
}

// Load reads the index from a file
func Load(path string) (i *Index, err error) {
	i = &Index{}
	xmlFile, err := os.Open(path)
	if err != nil {
		return
	}
	defer xmlFile.Close()
	dec := xml.NewDecoder(xmlFile)
	err = dec.Decode(i)
	return
}

// hashFile creates a sha1sum for a given file
func hashFile(path string) error {
	iFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer iFile.Close()
	h := sha1.New()
	_, err = io.Copy(h, iFile)
	if err != nil {
		return err
	}
	oFile, err := os.Create(path + ".sha1sum")
	if err != nil {
		return err
	}
	defer oFile.Close()
	fmt.Fprintf(oFile, "%x", h.Sum(nil))
	return nil
}

// Save writes the index out to a file, compresses it, and then generates hash files for both files
func (i *Index) Save(path string) error {
	indexFile := filepath.Join(path, "eopkg-index.xml")
	xmlFile, err := os.Create(indexFile)
	if err != nil {
		return err
	}
	enc := xml.NewEncoder(xmlFile)
	enc.Indent("    ", "    ")
	if err = enc.Encode(i); err != nil {
		xmlFile.Close()
		return err
	}
	xmlFile.Close()
	if err = shared.XzFile(indexFile, true); err != nil {
		return err
	}
	if err = hashFile(indexFile); err != nil {
		return err
	}
	return hashFile(indexFile + ".xz")
}
