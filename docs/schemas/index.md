# eopkg index XML

## Top-Level: PISI

```
<PISI>
    <Distribution>
        ... one only
    <Package>
        ... one or more
    <Component>
        ... one or more
    <Group>
        ... one or more
```

## Distribution

``` XML
<Distribution>
    <SourceName>Solus
    <Description>
        ... one or more
        @lang
        CDATA: Solus Repository
    <Version>1
    <Type>main
    <BinaryName>Solus
    <Obsoletes>
```

### Obsoletes

``` XML
<Obsoletes>
    <Package>pcre
        ... zero or more
```

## Package

``` XML
<Package>
    <Name>0ad
    <Summary>
        @xml:lang
        CDATA: A.D. is a free, open-source, cross-platform real-time strategy game
    <Description>
        @xml:lang
        CDATA:
    <PartOf>game.strategy
        ... name of a <Component>
    <License> BSD-2-Clause
        ... one or more
    <RuntimeDependencies>
    <History>
    <BuildHost>solus
    <Distribution>Solus
    <DistributionRelease>1
    <Architecture>x86_64
    <InstalledSize>9755641
    <PackageSize>2812298
    <PackageHash>37c9f2159869d8da09786d3646c5c26e8e8204d3
        SHA1SUM
    <PackageURI>0/0ad/0ad-0.0.23b-27-1-x86_64.eopkg
    <DeltaPackages>
    <PackageFormat>1.2
    <Source>
        <Name>0ad
        <Packager>
            <Name>Pierre-Yves
            <Email>
```

### RuntimeDependencies

``` XML
<RuntimeDependencies>
    <Dependency>
        ... zero or more
        @releaseFrom
            number of the release
        libsodium
```

### History

``` XML
<History>
    <Update>
        ... one or more
        @release
            number of the release
        <Date>YYYY-MM-DD
        <Version>0.0.23b
        <Comment>
            CDATA: Rebuild 0ad for libboost 1.72.0
        <Name> Pierre-Yves
        <Email>
```

### DeltaPackages

```XML
<DeltaPackages>
    <Delta>
        ... zero or more
        @releaseFrom 25
            number of the release
        <PackageURI>0/0ad/0ad-25-27-1-x86_64.delta.eopkg
        <PackageSize>2303734
        <PackageHash>d0e61c3e043d9550a382af7a3a4146ce2e58c8de
            SHA1SUM
```

## Component

``` XML
<Component>
    <Name>database
    <LocalName>
        ... one or more
        @xml:lang en
        CDATA: Database clients and servers
    <Summary>
        ... one or more
        @xml:lang en
        CDATA: Database clients and servers
    <Description>
        ... one or more
        @xml:lang en
        CDATA: Database clients and servers
    <Group>system
        Name of Group
    <Maintainer>
        <Name>Solus Team
        <Email>root@solus-project.com
```

## Groups

``` XML
<Group>
    <Name>security
    <LocalName>
        ... one or more
        @xml:lang en
        CDATA: Security Software
    <Icon>security-high
        entry somewhere in /usr/share/icons
```
