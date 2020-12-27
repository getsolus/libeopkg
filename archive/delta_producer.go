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
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DeltaProducer is responsible for taking two eopkg packages and spitting out
// a delta package for them, containing only the new files.
type DeltaProducer struct {
	left    *Package
	right   *Package
	baseDir string
	diffMap map[string]int
}

var (
	// ErrMismatchedDelta is returned when the input packages should never be delta'd,
	// i.e. they're unrelated
	ErrMismatchedDelta = errors.New("Delta is not possible between the input packages")

	// ErrDeltaPointless is returned when it is quite literally pointless to bother making
	// a delta package, due to the packages having exactly the same content.
	ErrDeltaPointless = errors.New("File set is the same, no point in creating delta")
)

// NewDeltaProducer will return a new delta producer for the given input packages
// It is very important that the old and new packages are in the correct order!
func NewDeltaProducer(baseDir string, left string, right string) (dp *DeltaProducer, err error) {
	var dirName string
	// Init a new DeltaProducer
	dp = &DeltaProducer{
		diffMap: make(map[string]int),
	}
	// Open the previous release
	dp.left, err = Open(left)
	if err != nil {
		goto CLOSE
	}
	// Read its contents
	if err = dp.left.ReadAll(); err != nil {
		goto CLOSE
	}
	// Open the new release
	dp.right, err = Open(right)
	if err != nil {
		goto CLOSE
	}
	// Read its contents
	if err = dp.right.ReadAll(); err != nil {
		goto CLOSE
	}
	// Check if these packages are from the same source
	if !IsDeltaPossible(&dp.left.Meta.Package, &dp.right.Meta.Package) {
		err = ErrMismatchedDelta
		goto CLOSE
	}
	// Form a unique directory entry
	dirName = fmt.Sprintf("%s-%s-%s-%d-%d",
		dp.left.Meta.Package.Name,
		dp.left.Meta.Package.GetVersion(),
		dp.left.Meta.Package.Architecture,
		dp.right.Meta.Package.GetRelease(),
		dp.left.Meta.Package.GetRelease())

	dp.baseDir = filepath.Join(baseDir, dirName)

	// Make sure base directory actually exists
	err = os.MkdirAll(dp.baseDir, 00755)

CLOSE:
	if err != nil {
		_ = dp.Close()
		dp = nil
	}
	return
}

// Close the DeltaProducer
func (d *DeltaProducer) Close() error {
	// check for nil pointer
	if d == nil {
		return nil
	}
	// close left
	if d.left != nil {
		d.left.Close()
		d.left = nil
	}
	// clode right
	if d.right != nil {
		d.right.Close()
		d.right = nil
	}
	// Ensure we always nuke the work directory we used
	if d.baseDir != "" {
		return os.RemoveAll(d.baseDir)
	}
	return nil
}

// produceTarball copies all of the files into the tarball and compresses it
func (d *DeltaProducer) produceTarball() (filename string, err error) {
	// Convert file lists to maps
	hashOldFiles := d.left.filesToMap()
	hashNewFiles := d.right.filesToMap()

	// Note this is very simple and works just like the existing eopkg functionality
	// which is purely hash-diff based. eopkg will look for relocations on applying
	// the update so that files get "reused"
	//
	// Special Note: Key "" denotes a directory which is basically empty, so we must
	// always include these in the delta
	for h, s := range hashNewFiles {
		// Skip identical files
		if _, ok := hashOldFiles[h]; ok && h != "" {
			continue
		}
		// Keep track of new files
		for _, p := range s {
			d.diffMap[strings.TrimSuffix(p.Path, "/")] = 1
		}
	}
	var installTar string
	var outF *os.File
	var tw *tar.Writer
	// All the same files
	if len(d.diffMap) == len(d.right.Files.File) {
		err = ErrDeltaPointless
		goto CLOSE
	}
	// No install.tar.xz to write as we have no different files
	if len(d.diffMap) == 0 {
		err = ErrDeltaPointless
		goto CLOSE
	}
	// Open output file to write our tarfile.
	installTar = filepath.Join(d.baseDir, "delta-eopkg.install.tar")
	outF, err = os.Create(installTar)
	if err != nil {
		goto CLOSE
	}
	tw = tar.NewWriter(outF)
	// Copy the delta files
	if err = d.copyModified(tw); err != nil {
		goto CLOSE
	}
	// Save the contents and close file
	tw.Flush()
	tw.Close()
	// Compress the tarball
	if err = XzFile(installTar, false); err != nil {
		goto CLOSE
	}
	filename = fmt.Sprintf("%s.xz", installTar)
	// Make sure we clean up properly!
CLOSE:
	if tw != nil {
		tw.Close()
	}
	return
}

// copyModified will iterate over the contents of the existing install.tar.xz
// for the new package, and only include the files that aren't hash-matched in the
// old files.xml
func (d *DeltaProducer) copyModified(tw *tar.Writer) error {
	// Unpack the existing tarball from zip container and decompressing it
	if err := d.right.ExtractTarball(d.baseDir); err != nil {
		return err
	}
	// Open the tarball
	inpFile := filepath.Join(d.baseDir, "install.tar")
	fi, err := os.Open(inpFile)
	if err != nil {
		return err
	}
	defer fi.Close()
	tarfile := tar.NewReader(fi)
	// Iterate over tarball contents
	for {
		// Get the next file
		header, err := tarfile.Next()
		if err != nil {
			// Check for no more files
			if err == io.EOF {
				err = nil
				break
			}
			// unexpected error
			return err
		}
		// Ensure that we compare things in the same way
		checkName := strings.TrimSuffix(header.Name, "/")
		// Skip anything not in the diff map
		if _, ok := d.diffMap[checkName]; !ok {
			continue
		}
		// Create new entry in the output tarball
		if err = tw.WriteHeader(header); err != nil {
			return err
		}
		// Copy file from package to delta
		if header.Typeflag == tar.TypeReg || header.Typeflag == tar.TypeRegA {
			if _, err = io.Copy(tw, tarfile); err != nil {
				return err
			}
		}
	}
	// flush contents to disk
	return tw.Flush()
}

// copyZipModified will iterate the central zip directory and skip only the
// install.tar.xz files, whilst copying everything else into the new zip
func (d *DeltaProducer) copyZipModified(zw *zip.Writer) error {
	// iterate over zip contents
	for _, zipFile := range d.right.zipFile.File {
		// Skip any kind of install.tar internally
		if strings.HasPrefix(zipFile.Name, "install.tar") {
			continue
		}
		// Open the file
		iop, err := zipFile.Open()
		if err != nil {
			return err
		}
		// Duplicate the header
		zwh := &zip.FileHeader{}
		*zwh = zipFile.FileHeader
		w, err := zw.CreateHeader(zwh)
		if err != nil {
			iop.Close()
			return err
		}
		// Copy the File member across (it implements FileHeader)
		if _, err = io.Copy(w, iop); err != nil {
			iop.Close()
			return err
		}
		// Close the File
		iop.Close() // be really sure we close it..
	}
	// flush to disk
	return zw.Flush()
}

// saveTarball adds the new tarball to the package
func (d *DeltaProducer) saveTarball(zipFile *zip.Writer, xzPath string) error {
	// No install.tar.xz to write as we have no different files
	if len(d.diffMap) == 0 {
		return ErrDeltaPointless
	}
	// Open the tarball file
	f, err := os.Open(xzPath)
	if err != nil {
		return err
	}
	defer f.Close()
	// Stat for header
	fst, err := f.Stat()
	if err != nil {
		return err
	}
	// Create a Zip file header
	fh, err := zip.FileInfoHeader(fst)
	if err != nil {
		return err
	}
	// Ensure it's always the right name.
	fh.Name = "install.tar.xz"
	// Write the header
	w, err := zipFile.CreateHeader(fh)
	if err != nil {
		return err
	}
	// Copy the tarball contents
	if _, err = io.Copy(w, f); err != nil {
		return err
	}
	// flush to disk
	return zipFile.Flush()
}

// Create will attempt to produce a delta between the 2 eopkg files
// This will be performed in temporary storage so must then be copied into
// the final resting location, and unlinked, before it can be used.
func (d *DeltaProducer) Create() (filename string, err error) {
	var out *os.File
	var zipFile *zip.Writer
	var zipPath string
	// Generate the tarball
	xzPath, err := d.produceTarball()
	if err != nil {
		goto CLOSE
	}
	// Open the Zip file
	filename = ComputeDeltaName(&d.left.Meta.Package, &d.right.Meta.Package)
	zipPath = filepath.Join(d.baseDir, filename)
	filename = zipPath
	out, err = os.Create(zipPath)
	if err != nil {
		goto CLOSE
	}
	zipFile = zip.NewWriter(out)
	// Copy the other files
	if err = d.copyZipModified(zipFile); err != nil {
		goto CLOSE
	}
	// Copy the tarball
	if err = d.saveTarball(zipFile, xzPath); err != nil {
		goto CLOSE
	}
	// Close the zip file
	err = zipFile.Close()

CLOSE:
	// If we're successful, we don't delete these
	if err != nil {
		_ = os.Remove(xzPath)
		_ = os.Remove(zipPath)
	}
	return
}
