# TODO

## Archive

### General

 - [x] Refactor modePrivate in File so use custom Marshal/Unmarshal instead of private field
 - [x] New Diff for Files when generating Deltas

### Required for ferryd

 - [x] Extract `files.xml` and `metadata.xml`
 - [x] Generate Deltas
    - [x] Calculate difference between two packages
    - [x] Copy existing `files.xml` from new package
    - [x] Copy existing `metadata.xml` from new package
    - [x] Build a package from differences
 - [ ] Write, compress, and sha256sum Index

### Required for ypkg3

 - [ ] Build a package
    - [ ] Generate and compress tarball from the contents of a directory
    - [ ] Generating `files.xml` from the contents of a directory
    - [ ] Writing Metadata to `metadata.xml`
    - [ ] Updating `pspec_x86_64.xml`

### Required for sol

 - [ ] Read in the the index
 - [ ] Installation
    - [ ] Unpack tar to directory
    - [ ] Verify files during/after writing to disk
    - [ ] Set file characteristics from `files.xml`
 - [ ] Upgrades
    - [ ] Handle Delta Installation
 - [ ] Removals
    - [ ] Calculating Deletions from `files.xml`
