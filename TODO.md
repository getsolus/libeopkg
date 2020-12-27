# TODO

## Archive

### General

 - [ ] Refactor modePrivate in File so use custom Marshal/Unmarshal instead of private field
 - [ ] Checksum archive file


### Required for ferryd

 - [ ] Extract and compress: `files.xml` and `metadata.xml`
 - [ ] Generate Deltas
    - [ ] Calculate difference between two packages
    - [ ] Generate new `files.xml` from diff
    - [ ] Copy existing `metadata.xml` from new package
    - [ ] Build a package from differences
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
