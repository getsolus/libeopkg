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
	"encoding/xml"
	"github.com/getsolus/libeopkg/shared"
)

// Metadata contains all of the information a package can provide to a user
// prior to installation. This includes the name, version, release, and so
// forth.
//
// Every Package contains Metadata, and during eopkg indexing, a reduced
// version of the Metadata is emitted.
type Metadata struct {
	Source  shared.Source
	Package *Package `xml:"Package"`
}

// ReadMetadata will read the `metadata.xml` file within the archive and
// deserialize it into something accessible within the .eopkg container.
func (a *Archive) ReadMetadata() error {
	// Open the metadata file
	metaFile := a.FindFile("metadata.xml")
	if metaFile == nil {
		return shared.ErrEopkgCorrupted
	}
	f, err := metaFile.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	// Decode its contents
	a.Meta = &Metadata{}
	dec := xml.NewDecoder(f)
	if err = dec.Decode(a.Meta); err != nil {
		return err
	}
	// Remove extraneous spaces and fix missing localised fields
	a.Meta.Package.Clean()
	return nil
}
