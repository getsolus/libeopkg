# TODO

## General

 - [x] Refactor modePrivate in File so use custom Marshal/Unmarshal instead of private field
 - [x] New Diff for Files when generating Deltas
 - [x] Refresh documentation after fighting with the XML parser

## Required for ferryd

 - [x] Extract `files.xml` and `metadata.xml`
 - [x] Generate Deltas
    - [x] Calculate difference between two packages
    - [x] Copy existing `files.xml` from new package
    - [x] Copy existing `metadata.xml` from new package
    - [x] Build a package from differences
 - [x] Write, compress, and sha1sum Index

## Required for ypkg3

 - [ ] Build a package
    - [ ] Generate and compress tarball from the contents of a directory
    - [ ] Generating `files.xml` from the contents of a directory
    - [ ] Writing Metadata to `metadata.xml`
    - [ ] Updating `pspec_x86_64.xml`

## Required for sol

 - [x] Read in the the index
 - [x] Installation
    - [x] Unpack tar to directory
    - [x] Verify files during/after writing to disk (now with symlinks)
    - [x] Set file characteristics from `files.xml` (skip because TAR is better)
    - [x] Install metadata files to different directory
 - [ ] Upgrades
    - [ ] Handle Delta Installation
 - [ ] Removals
    - [ ] Calculating Deletions from `files.xml`
 - [x] Verify files after installation (eopkg check)
