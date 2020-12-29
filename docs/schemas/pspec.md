# pspec.*.xml

## Filename

1. `pspec.xml` (legacy format for PISI aka eopkg build, NOT documented here)
2. `pspec_x86_64.xml` (produced by ypkg after the build)

## Format

``` XML
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
    <History>
```

---

### Source

``` XML
<Source>
	<Name> Name of the Source Repo / Package?
	<Homepage> URL (optional)
	<Packager>
		<Name> Maintainer
		<Email>

```

---

### Package

``` XML
<Package>
	<Name> nano
	<Summary xml:lang="en"> Small, friendly text editor inspired by Pico
	<Description xml:lang="en"> GNU nano is an easy-to-use text editor...
	<PartOf> system.devel
	<License> GPL-3.0-or-later
	<RuntimeDependencies>
		... omit if empty
		<Dependency release="62"> glibc
			... one or more
	<Files>
```

#### Files

```XML
<Files>
	... one or more
	<Path>/usr/bin/nano
		@fileType "executable"
			kind of file based on location in tree
```

**FileTypes**

- **config**
	+ /etc
- **data**
	+ /usr/lib/pkgconfig
	+ /usr/lib32/pkgconfig
	+ /usr/lib64/pkgconfig
- **doc**
	+ /usr/share/doc
	+ /usr/share/gtk-doc
	+ /usr/share/help
- **executable**
	* /bin
	* /sbin
	+ /usr/bin
	+ /usr/libexec
	+ /usr/sbin
- **header**
	+ /usr/include
- **info**
	+ /usr/share/info
- **library**
	+ /usr/lib
- **localedata**
	+ /usr/share/locale
- **man**
	+ /usr/share/man
	
**NOTE**: Unless otherwise listed, all remaining files will be marked as `data`.

---

### History

``` XML
<History>
	<Update>
		... zero or more
		@release 123
		<Date> YYYY-MM-DD
		<Version> 5.3
		<Comment> Packaging update
		<Name> Maintainer
		<Email>
```

History can take one of two forms in `pspec_x86_64.xml`:

1. A single entry for the new release that will be created
2. The complete history, in the case that someone forgets to increment the release number.