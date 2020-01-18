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
	"archive/zip"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//
// A Package is used for accessing a `.eopkg` archive, the current format used
// within Solus for software packages.
//
// An .eopkg archive is actually a ZIP archive. Internally it has the following
// structure:
//
//      metadata.xml    -> Package information
//      files.xml       -> Record of the files and hash/uid/gid/etc
//      comar/          -> Postinstall scripts
//      install.tar.xz  -> Filesystem contents
//
// Due to this toplevel simplicity, we can use golang's native `archive/zip`
// library to achieve eopkg access, and parse the contents accordingly.
// This is much faster than having to call out to the host side tool, which
// is presently written in Python.
//
type Package struct {
	// Path to this .eopkg file
	Path string
	// Basename of the package, unique.
	ID string
	// Metadata for this package
	Meta *Metadata
	// Files for this package
	Files *Files
	// .eopkg is a zip archive
	zipFile *zip.ReadCloser
}

// Open will attempt to open the given .eopkg file.
// This must be a valid .eopkg file and this stage will assert that it is
// indeed a real archive.
func Open(path string) (*Package, error) {
	// Create package object
	ret := &Package{
		Path: path,
		ID:   filepath.Base(path),
	}
	// Open package file
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	ret.zipFile = zipFile
	return ret, nil
}

// Close a previously opened .eopkg file
func (p *Package) Close() error {
	return p.zipFile.Close()
}

// FindFile will search for the given name in the .zip's
// file headers.
// We do not need to worry about the issue with the Name
// member being the basename, as the filenames are always
// unique.
//
// In the event of the file requested not being found,
// we return nil. The caller should then bail and indicate
// that the eopkg is corrupted.
func (p *Package) FindFile(path string) *zip.File {
	// Iterate over all files
	for _, f := range p.zipFile.File {
		// Check for match
		if path == f.Name {
			return f
		}
	}
	return nil
}

// ReadMetadata will read the `metadata.xml` file within the archive and
// deserialize it into something accessible within the .eopkg container.
func (p *Package) ReadMetadata() error {
	// Already read metadata
	if p.Meta != nil {
		return nil
	}
	// Open the metadata file
	metaFile := p.FindFile("metadata.xml")
	if metaFile == nil {
		return ErrEopkgCorrupted
	}
	fi, err := metaFile.Open()
	if err != nil {
		return err
	}
	defer fi.Close()
	// Decode its contents
	metadata := &Metadata{}
	dec := xml.NewDecoder(fi)
	if err = dec.Decode(metadata); err != nil {
		return err
	}
	p.Meta = metadata
	// Remove extraneous spaces and fix missing localised fields
	for i := range p.Meta.Package.Summary {
		sum := &p.Meta.Package.Summary[i]
		sum.Value = strings.TrimSpace(sum.Value)
	}
	FixMissingLocalLanguage(&p.Meta.Package.Summary)
	for i := range p.Meta.Package.Description {
		desc := &p.Meta.Package.Description[i]
		desc.Value = strings.TrimSpace(desc.Value)
	}
	FixMissingLocalLanguage(&p.Meta.Package.Description)
	return nil
}

// ReadFiles will read the `files.xml` file within the archive and
// deserialize it into something accessible within the .eopkg container.
func (p *Package) ReadFiles() error {
	// Already read Files
	if p.Files != nil {
		return nil
	}
	// Open the files list
	files := p.FindFile("files.xml")
	if files == nil {
		return ErrEopkgCorrupted
	}
	fi, err := files.Open()
	if err != nil {
		return err
	}
	defer fi.Close()
	// Decode its contents
	ret := &Files{}
	dec := xml.NewDecoder(fi)
	if err = dec.Decode(ret); err != nil {
		return err
	}
	// Convert file modes from strings to ints
	for _, f := range ret.File {
		if err := f.ParseFileMode(); err != nil {
			return err
		}
	}
	p.Files = ret
	return nil
}

// ReadAll will read both the metadata + files xml files
func (p *Package) ReadAll() error {
	if err := p.ReadMetadata(); err != nil {
		return err
	}
	return p.ReadFiles()
}

// ExtractTarball will fully extract install.tar.xz to the destination
// direction + install.tar suffix
func (p *Package) ExtractTarball(directory string) error {
	// Open tarball
	xzName := filepath.Join(directory, "install.tar.xz")
	tarball := p.FindFile("install.tar.xz")
	if tarball == nil {
		return ErrEopkgCorrupted
	}

	fi, err := tarball.Open()
	if err != nil {
		return err
	}
	defer fi.Close()
	// Create destination tarball
	outF, err := os.Create(xzName)
	if err != nil {
		return err
	}
	defer outF.Close()
	// Copy the entire tarball
	if _, err = io.Copy(outF, fi); err != nil {
		return err
	}
	// Uncompress the tarball
	return UnxzFile(xzName, false)
}

// filesToMap is a helper that will let us uniquely index hash to file-set
func (p *Package) filesToMap() (ret map[string][]*File) {
	ret = make(map[string][]*File)
	// For each file in files
	for _, f := range p.Files.File {
		// Append it to a list of files with the same hash
		ret[f.Hash] = append(ret[f.Hash], f)
	}
	return ret
}
