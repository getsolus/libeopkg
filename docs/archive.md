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
./comar
./comar/manager.py
./comar/package.py
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
            One of: see `pspec.xml`
        <Size>320264
        <Uid>0
        <Gid>0
        <Mode>0755
        <Hash>ffe45818d3fdfe4ee01efb6473176513d0ace947
            SHA1SUM
        <Permanent>
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
	<Package>
```

## Source

``` XML
<Source>
    <Name>nano
    <Homepage>https://www.nano-editor.org
    <Packager>
        <Name>Rune Morling
        <Email>
```

## Package

``` XML
<Package>
    <Name> pisi
    <Summary>
    		@xml:lang
    		CDATA: PISI
    <Description>
    		@xml:lang
    		CDATA: PISI is a modern package management...
		<IsA> app:console (legacy)
		<PartOf>system.base
		<License>
			... one or more
			GPL-2.0-or-later
	<RuntimeDependencies>
    <Replaces>
        ... omit if empty
        ... one or more packages
    		<Package>
    <Conflicts>
	    ... omit if empty
	    ... one or more packages
	    <Package>
	<Provides>
	<History>
	<BuildHost> solus-build-server
	<Distribution> Solus
	<DistributionRelease> 1
	<Architecture> x86_64
	<InstalledSize> 1709169
	<PackageFormat> 1.2
	<Source>
```

### RuntimeDependencies

``` XML
<RuntimeDependencies>
	... omit if empty
	<Dependency>
		... one or more
		@releaseFrom 13
		ncurses
```

### Provides

``` XML
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

### History

``` XML
<History>
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
```

