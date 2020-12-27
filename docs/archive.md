# .eopkg Archive File Format

# Filenames
``` XML
<name>-<version>-<release>-<distribution-release>-<architecture>.eopkg
<name>-<releaseFrom>-<release>-<distribution-release>-<architecture>.delta.eopkg
```

# Delta Packages
Identical to a normal eopkg, except they only contain the files that have been modified.

# Zip Archive

```
./files.xml
./install.tar.xz
./metadata.xml
```

# files.xml

``` XML
<Files>
    <File>
        ... zero or more
        <Path>usr/bin/nano
        <Type>executable
            One of: executable, doc, data, info, localedata
        <Size>320264
        <Uid>0
        <Gid>0
        <Mode>0755
        <Hash>ffe45818d3fdfe4ee01efb6473176513d0ace947
            SHA1SUM
        <Permananent>
            Do not remove on update, maybe deletion
```

# install.tar.*

- Files and directories, relative to `/`
- Compressed with any one of a nubmer of formats
    1. Gzip
    2. Bzip2
    3. XZ (post-pisi)

# metadata.xml

```
<PISI>
    <Source>
        <Name>nano
        <Homepage>https://www.nano-editor.org
        <Packager>
            <Name>Rune Morling
            <Email>
    <Package>
        <Dependency>
            @releaseFrom
            @releaseTo
            @release
            @versionFrom
            @versionTo
            @version
        <Update>
            @release int
            @type string
            <Date>
            <Version>
            <Comment> CDATA: commit message
            <Name> CDATA: maintainer
            <Email> maintainer email
            <Requires>
                <Action>
                    @package name of package
                    CDATA: value
        <Provides>
            <COMAR>
                ... zero or more
                @script name of the script to run
                CDATA:
            <PkgConfig> pkgconfig(glib-2.0)
                ... zero or more
            <PkgConfig32>
                ... zero or more
```
