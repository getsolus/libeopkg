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
	"bytes"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/getsolus/libeopkg/shared"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

// File is the idoimatic representation of the XML <File> node
//
// Note that directories are indicated by a missing hash. Unfortunately
// however eopkg doesn't record the actual _type_ of a file in an intelligent
// sense, thus we'll have to deal with symlinks separately.
//
// In an ideal world the package archive would be hash indexed with no file
// names or special permissions inside the archive, and we'd record all relevant
// metadata. This would allow a single copy, many hardlink approach to blit
// the files out, as well as allowing us to more accurately represent symbolic
// links instead of pretending they're real files.
//
// Long story short: Wait for eopkg's successor to worry about this stuff.
type File struct {
	Path      string
	Type      shared.FileType
	Size      int64    `xml:",omitempty"`
	UID       int      `xml:"UID,omitempty"`
	GID       int      `xml:"GID,omitempty"`
	Mode      FileMode `xml:",omitempty"`
	Hash      string   `xml:",omitempty"`
	Permanent string   `xml:",omitempty"`
}

// Equal checks if one file is identical to another
func (f *File) Equal(other *File) bool {
	return f.Path == other.Path && f.Type == other.Type && f.Size == other.Size &&
		f.UID == other.UID && other.GID == other.GID && f.Mode == other.Mode &&
		f.Hash == other.Hash && f.Permanent == other.Permanent
}

// FileMode will return an os.FileMode version of our string encoded "Mode" member
func (f *File) FileMode() os.FileMode {
	return os.FileMode(f.Mode)
}

// Verify reads the copy of this file from disk, hashes it, and compares for the correct hash
func (f *File) Verify(path string) error {
	dstPath := filepath.Join(path, f.Path)
	info, err := os.Lstat(dstPath)
	if err != nil {
		return err
	}
	if sz := info.Size(); sz != f.Size {
		return fmt.Errorf("'%s' size mismatch: %d != %d", f.Path, sz, f.Size)
	}
	if m := info.Mode().Perm(); m != f.FileMode() {
		return fmt.Errorf("'%s' permission mismatch: %d != %d", f.Path, m, f.Mode)
	}
	stat := info.Sys().(*syscall.Stat_t)
	if stat.Uid != uint32(f.UID) {
		return fmt.Errorf("'%s' UID mismatch: %d != %d", f.Path, stat.Uid, uint32(f.UID))
	}
	if stat.Gid != uint32(f.GID) {
		return fmt.Errorf("'%s' GID mismatch: %d != %d", f.Path, stat.Gid, uint32(f.GID))
	}
	var in io.Reader
	if mode := info.Mode(); !mode.IsRegular() {
		if (mode & os.ModeSymlink) != os.ModeSymlink {
			return nil
		}
		name, err := os.Readlink(dstPath)
		if err != nil {
			return err
		}
		in = bytes.NewBuffer([]byte(name))
	} else {
		dst, err := os.Open(dstPath)
		if err != nil {
			return err
		}
		defer dst.Close()
		in = dst
	}
	h := sha1.New()
	if _, err := io.Copy(h, in); err != nil {
		return err
	}
	sum := fmt.Sprintf("%x", h.Sum(nil))
	if sum != f.Hash {
		return fmt.Errorf("'%s' hash mismatch: %s != %s", f.Path, sum, f.Hash)
	}
	return nil
}

// ReadFiles will read the `files.xml` file within the archive and
// deserialize it into something accessible within the .eopkg container.
func (p *Archive) ReadFiles() error {
	// Already read Files
	if p.Files != nil {
		return nil
	}
	// Open the files list
	files := p.FindFile("files.xml")
	if files == nil {
		return shared.ErrEopkgCorrupted
	}
	f, err := files.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	// Decode its contents
	p.Files = &Files{}
	dec := xml.NewDecoder(f)
	return dec.Decode(p.Files)
}

// Files is the idiomatic representation of the XML <Files> node with one or
// more <File> children
type Files struct {
	File []*File
}

// HasFile checks if the specified path is listed
func (fs Files) HasFile(path string) bool {
	for _, f := range fs.File {
		if f.Path == path {
			return true
		}
	}
	return false
}

// Diff creates a new Files from all of the modifications between "other" and this Files
func (fs *Files) Diff(other *Files) (modified, removed *Files) {
	modified, removed = &Files{}, &Files{}
	// Check for modified or removed files
	for _, curr := range fs.File {
		found := false
		for _, next := range other.File {
			if curr.Path != next.Path {
				continue
			}
			if !curr.Equal(next) {
				modified.File = append(modified.File, next)
				found = true
				break
			}
		}
		if !found {
			removed.File = append(removed.File, curr)
		}
	}
	// Check for new files
	for _, next := range other.File {
		found := false
		for _, curr := range fs.File {
			if next.Equal(curr) {
				found = true
				break
			}
		}
		if !found {
			modified.File = append(modified.File, next)
		}
	}
	return
}

// FileMode is a hexdecimal excoded FileMode
type FileMode os.FileMode

// MarshalXML allows writing a FileMode as Hex
func (fm FileMode) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	raw := strconv.FormatUint(uint64(fm), 8)
	if len(raw) == 3 || len(raw) == 4 {
		raw = "0" + raw
	}
	return e.EncodeElement(raw, start)
}

// UnmarshalXML allows reading a FileMode as Hex
func (fm *FileMode) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error
	var raw string
	if err = d.DecodeElement(&raw, &start); err != nil {
		return err
	}
	var i uint64
	if i, err = strconv.ParseUint(raw, 8, 16); err != nil {
		return err
	}
	*fm = FileMode(i)
	return nil
}

func (fm FileMode) String() string {
	return os.FileMode(fm).String()
}
