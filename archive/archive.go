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
	"archive/tar"
	"archive/zip"
	"fmt"
	"github.com/getsolus/libeopkg/shared"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
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
	// Unpack tarball
	if err := a.UnpackFile("install.tar.xz", directory); err != nil {
		return err
	}
	// Uncompress the tarball
	xzName := filepath.Join(directory, "install.tar.xz")
	return shared.UnxzFile(xzName, false)
}

// setXattrs, applies Xattrs to a file, as they are stored in `tar.Header.PAXRecords`
func setXattrs(path string, paxattrs map[string]string) error {
	for attr, value := range paxattrs {
		fmt.Printf("attr: %s, value: %s\n", attr, value)
		if !strings.HasPrefix(attr, "SCHILY.xattr.") {
			continue
		}
		name := strings.TrimPrefix(attr, "SCHILY.xattr.")
		if err := syscall.Setxattr(path, name, []byte(value), 0); err != nil {
			return err
		}
	}
	return nil
}

// UnpackFile copies file from Zip to destination
func (a *Archive) UnpackFile(name, path string) error {
	srcFile := a.FindFile(name)
	if srcFile == nil {
		return shared.ErrEopkgCorrupted
	}
	src, err := srcFile.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dstFile := filepath.Join(path, name)
	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return dst.Sync()
}

// Unpack writes out 'files.xml' and 'metadata.xml', then unpacks the tarball to "install"
func (a *Archive) Unpack(metaPath, filesPath string) error {
	// Unpack metadata files
	// Make dir to unpack things into
	if err := os.MkdirAll(metaPath, 0755); err != nil {
		return err
	}
	if err := a.UnpackFile("files.xml", metaPath); err != nil {
		return err
	}
	if err := a.UnpackFile("metadata.xml", metaPath); err != nil {
		return err
	}
	// Make subdir to unpack things into
	if err := os.MkdirAll(filesPath, 0755); err != nil {
		return err
	}
	// Unpack the existing tarball from zip container and decompressing it
	if err := a.ExtractTarball(filesPath); err != nil {
		return err
	}
	// Open the tarball
	srcFile := filepath.Join(filesPath, "install.tar")
	f, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer os.Remove(srcFile)
	defer f.Close()
	src := tar.NewReader(f)
	// Iterate over tarball contents
	header, err := src.Next()
	for err == nil && header != nil {
		dstPath := filepath.Join(filesPath, header.Name)
		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}
		// make the output directory
		switch header.Typeflag {
		case tar.TypeReg, tar.TypeRegA:
			dst, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			_, err = io.Copy(dst, src)
			_ = dst.Sync()
			_ = dst.Close()
		case tar.TypeLink:
			err = os.Link(header.Linkname, dstPath)
		case tar.TypeSymlink:
			err = os.Symlink(header.Linkname, dstPath)
		case tar.TypeDir:
			err = os.MkdirAll(dstPath, 0755)
		case tar.TypeFifo:
			err = syscall.Mkfifo(dstPath, 0666)
		default:
			// unexpected
			continue
		}
		if err != nil {
			return err
		}
		if err = os.Chown(dstPath, header.Uid, header.Gid); err != nil {
			return err
		}
		if err = os.Chtimes(dstPath, header.AccessTime, header.ModTime); err != nil {
			return err
		}
		if err = setXattrs(dstPath, header.PAXRecords); err != nil {
			return err
		}
		header, err = src.Next()
	}
	if err == io.EOF {
		err = nil
	}
	return err
}

// Verify validates all of the files on disk against the archive
func (a *Archive) Verify(path string) error {
	if err := a.ReadAll(); err != nil {
		return err
	}
	for _, file := range a.Files.File {
		if err := file.Verify(path); err != nil {
			return err
		}
	}
	return nil
}

// IsDeltaPossible checks is a delta can be made from the two provided packages
func (a *Archive) IsDeltaPossible(newRelease *Archive) bool {
	return a.Meta.Package.IsDeltaPossible(newRelease.Meta.Package)
}

// Diff gets a list of files in "other" that have been modified or don't exist in this archive
func (a *Archive) Diff(other *Archive) (*Files, *Files) {
	return a.Files.Diff(other.Files)
}
