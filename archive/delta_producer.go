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

var (
	// ErrMismatchedDelta is returned when the input packages should never be delta'd,
	// i.e. they're unrelated
	ErrMismatchedDelta = errors.New("Delta is not possible between the input packages")

	// ErrDeltaPointless is returned when it is quite literally pointless to bother making
	// a delta package, due to the packages having exactly the same content.
	ErrDeltaPointless = errors.New("File set is the same, no point in creating delta")
)

// DeltaProducer is responsible for taking two eopkg packages and spitting out
// a delta package for them, containing only the new files.
type DeltaProducer struct {
	left    *Archive
	right   *Archive
	prefix  string
	workDir string
	diffMap map[string]int
}

// NewDeltaProducer will return a new delta producer for the given input packages
// It is very important that the old and new packages are in the correct order!
func NewDeltaProducer(workDir string, left string, right string) (dp *DeltaProducer, err error) {
	// Init a new DeltaProducer
	dp = &DeltaProducer{
		diffMap: make(map[string]int),
	}
	// Open the previous release
	dp.left, err = OpenAll(left)
	if err != nil {
		goto CLOSE
	}
	// Open the new release
	dp.right, err = OpenAll(right)
	if err != nil {
		goto CLOSE
	}
	// Check if these packages are from the same source
	if !dp.left.IsDeltaPossible(dp.right) {
		err = ErrMismatchedDelta
		goto CLOSE
	}
	// Form a unique directory entry
	dp.prefix = dp.left.Meta.Package.DeltaName(dp.right.Meta.Package.GetRelease())
	dp.workDir = filepath.Join(workDir, dp.prefix)

	// Make sure base directory actually exists
	err = os.MkdirAll(dp.workDir, 00755)

CLOSE:
	if err != nil {
		_ = dp.Close()
		dp = nil
	}
	return
}

// Close the DeltaProducer
func (dp *DeltaProducer) Close() error {
	// check for nil pointer
	if dp == nil {
		return nil
	}
	dp.left.Close()
	dp.right.Close()
	// Ensure we always nuke the work directory we used
	if dp.workDir != "" {
		return os.RemoveAll(dp.workDir)
	}
	return nil
}

// produceTarball copies all of the files into the tarball and compresses it
func (dp *DeltaProducer) produceTarball() (filename string, err error) {
	// Convert file lists to maps
	modified, removed := dp.left.Diff(dp.right)

	// All the same files
	if len(modified.File) == 0 && len(removed.File) == 0 {
		err = ErrDeltaPointless
		return
	}
	// Open output file to write our tarfile.
	installTar := filepath.Join(dp.workDir, "delta-install.tar")
	outF, err := os.Create(installTar)
	if err != nil {
		return
	}
	tw := tar.NewWriter(outF)
	// Copy the delta files
	if err = dp.Copy(tw, modified); err != nil {
		outF.Close()
		return
	}
	// Save the contents and close file
	tw.Close()
	outF.Close()
	// Compress the tarball
	filename = fmt.Sprintf("%s.xz", installTar)
	println(filename)
	err = XzFile(installTar, false)
	return
}

// Copy will iterate over the contents of the existing install.tar.xz for the new package,
// and only include the files that aren't hash-matched in the old files.xml
func (dp *DeltaProducer) Copy(dst *tar.Writer, modified *Files) error {
	// Unpack the existing tarball from zip container and decompressing it
	if err := dp.right.ExtractTarball(dp.workDir); err != nil {
		return err
	}
	// Open the tarball
	srcTar := filepath.Join(dp.workDir, "delta-install.tar")
	f, err := os.Open(srcTar)
	if err != nil {
		return err
	}
	defer f.Close()
	src := tar.NewReader(f)
	// Iterate over tarball contents
	header, err := src.Next()
	for err != nil && header != nil {
		// Ensure that we compare things in the same way
		path := strings.TrimSuffix(header.Name, "/")
		// Skip anything not in modifed
		if !modified.HasFile(path) {
			continue
		}
		// Create new entry in the output tarball
		if err = dst.WriteHeader(header); err != nil {
			return err
		}
		// Copy file from package to delta
		if header.Typeflag == tar.TypeReg || header.Typeflag == tar.TypeRegA {
			if _, err = io.Copy(dst, src); err != nil {
				return err
			}
		}
		header, err = src.Next()
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return err
	}
	// flush contents to disk
	return dst.Flush()
}

// copyZipModified will iterate the central zip directory and skip only the
// install.tar.xz files, whilst copying everything else into the new zip
func (dp *DeltaProducer) copyZipModified(dst *zip.Writer) error {
	// iterate over zip contents
	for _, src := range dp.right.zipFile.File {
		// Skip any kind of install.tar internally
		if strings.HasPrefix(src.Name, "install.tar") {
			continue
		}
		// Open the file
		srcFile, err := src.Open()
		if err != nil {
			return err
		}
		// Duplicate the header
		dstFile, err := dst.CreateHeader(&src.FileHeader)
		if err != nil {
			srcFile.Close()
			return err
		}
		// Copy the File member across (it implements FileHeader)
		if _, err = io.Copy(dstFile, srcFile); err != nil {
			srcFile.Close()
			return err
		}
		// Close the File
		srcFile.Close() // be really sure we close it..
	}
	// flush to disk
	return dst.Flush()
}

// saveTarball adds the new tarball to the package
func (dp *DeltaProducer) saveTarball(zipFile *zip.Writer, xzPath string) error {
	// Open the tarball file
	src, err := os.Open(xzPath)
	if err != nil {
		return err
	}
	defer src.Close()
	// Stat for header
	fst, err := src.Stat()
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
	dst, err := zipFile.CreateHeader(fh)
	if err != nil {
		return err
	}
	// Copy the tarball contents
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	// flush to disk
	return zipFile.Flush()
}

// Create will attempt to produce a delta between the 2 eopkg files
// This will be performed in temporary storage so must then be copied into
// the final resting location, and unlinked, before it can be used.
func (dp *DeltaProducer) Create() (filename string, err error) {
	// Generate the tarball
	xzPath, err := dp.produceTarball()
	if err != nil {
		_ = os.Remove(xzPath)
		return
	}
	// Open the Zip file
	filename = filepath.Join(dp.workDir, dp.prefix+".delta.eopkg")
	dst, err := os.Create(filename)
	if err != nil {
		dp.cleanup(xzPath, filename)
		return
	}
	defer dst.Close()
	zipFile := zip.NewWriter(dst)
	// Copy the other files
	if err = dp.copyZipModified(zipFile); err != nil {
		dp.cleanup(xzPath, filename)
		zipFile.Close()
		return
	}
	// Copy the tarball
	if err = dp.saveTarball(zipFile, xzPath); err != nil {
		dp.cleanup(xzPath, filename)
		zipFile.Close()
		return
	}
	// Close the zip file
	err = zipFile.Close()
	return
}

// cleanup removes imcomplete deta files
func (dp *DeltaProducer) cleanup(xzPath, zipPath string) {
	_ = os.Remove(xzPath)
	_ = os.Remove(zipPath)
}
