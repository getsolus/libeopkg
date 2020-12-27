# pspec.*.xml

## Filename

1. `pspec.xml` (legacy format for PISI aka eopkg build)
2. `pspec_x86_64.xml`

## Format

```
<PISI>
    <Source>
    <Package>
        Identical to Index **EXCEPT**
        1. No DeltaPackages
        2. Files are listed
        3. History only contains the latest entry UNLESS rebuilding against the same release
            - Then: Contains the entire history
    <Files>
        <Path>
            @fileType executable
                <File><Type> from files.xml
            /usr/bin/nano
```