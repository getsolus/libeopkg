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
	"archive/zip"
	"github.com/getsolus/libeopkg/shared"
	"io"
	"os"
	"path/filepath"
)

//
// An Archive is used for accessing a `.eopkg` archive, the current format used
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
type Archive struct {
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
func Open(path string) (a *Archive, err error) {
	// Create package object
	a = &Archive{
		Path: path,
		ID:   filepath.Base(path),
	}
	// Open package file
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return
	}
	a.zipFile = zipFile
	return
}

// OpenAll will Open an Archive and ReadAll of its metadata
func OpenAll(path string) (a *Archive, err error) {
	if a, err = Open(path); err != nil {
		return
	}
	err = a.ReadAll()
	return
}

// Close a previously opened .eopkg file
func (a *Archive) Close() error {
	if a == nil {
		return nil
	}
    if a.zipFile != nil {
        err := a.zipFile.Close()
        a.zipFile = nil
        return err
    }
	return nil
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
func (a *Archive) FindFile(path string) *zip.File {
	// Iterate over all files
	for _, f := range a.zipFile.File {
		// Check for match
		if path == f.Name {
			return f
		}
	}
	return nil
}

// ReadAll will read both the metadata + files xml files
func (a *Archive) ReadAll() error {
	if err := a.ReadMetadata(); err != nil {
		return err
	}
	return a.ReadFiles()
}

// ExtractTarball will fully extract install.tar.xz to the destination direction + install.tar suffix
func (a *Archive) ExtractTarball(directory string) error {
	// Open tarball
	tarball := a.FindFile("install.tar.xz")
	if tarball == nil {
		return shared.ErrEopkgCorrupted
	}
	f, err := tarball.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	// Create destination tarball
	xzName := filepath.Join(directory, "install.tar.xz")
	outF, err := os.Create(xzName)
	if err != nil {
		return err
	}
	defer outF.Close()
	// Copy the entire tarball
	if _, err = io.Copy(outF, f); err != nil {
		return err
	}
	// Uncompress the tarball
	return UnxzFile(xzName, false)
}

// IsDeltaPossible checks is a delta can be made from the two provided packages
func (a *Archive) IsDeltaPossible(newRelease *Archive) bool {
	return a.Meta.Package.IsDeltaPossible(newRelease.Meta.Package)
}

// Diff gets a list of files in "other" that have been modified or don't exist in this archive
func (a *Archive) Diff(other *Archive) (*Files, *Files) {
	return a.Files.Diff(other.Files)
}
