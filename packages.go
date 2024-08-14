package main

import (
	"fmt"
	"strings"
)

// Package is the NEVRA representation of an RPM package string
type Package struct {
	Name    string `json:"name"`
	Epoch   string `json:"epoch"`
	Version string `json:"version"`
	Release string `json:"release"`
	Arch    string `json:"arch"`
	Raw     string `json:"raw"`
}

// NewPackageFromString returns the parsed NEVRA from an RPM package string
func NewPackageFromString(pkgstr string) (pkg Package, err error) {
	pkg = Package{}
	pkg.Raw = pkgstr
	pkg.Arch, pkgstr, err = popArch(pkgstr)
	if err != nil {
		return
	}
	pkg.Release, pkgstr, err = popRelease(pkgstr)
	if err != nil {
		return
	}
	pkg.Version, pkgstr, err = popVersion(pkgstr)
	if err != nil {
		return
	}
	pkg.Epoch, pkg.Name = parseEpochName(pkgstr)
	return
}

// sameRelease returns whether two packages are from the same amazonlinux release.
func sameRelease(pkgA Package, pkgB Package) bool {
	if strings.Contains(pkgA.Release, "amzn2023") && strings.Contains(pkgB.Release, "amzn2023") {
		return true
	}
	if strings.Contains(pkgA.Release, "amzn2023") && strings.Contains(pkgB.Release, "amzn2022") {
		return true
	}
	if strings.Contains(pkgA.Release, "amzn2") && strings.Contains(pkgB.Release, "amzn2") {
		return true
	}
	if strings.Contains(pkgA.Release, "amzn1") && strings.Contains(pkgB.Release, "amzn1") {
		return true
	}
	return false
}

// Pops the architecture from the end of a full RPM string
func popArch(raw string) (popped string, arch string, err error) {
	popped, arch, err = popDelim(raw, ".")
	if err != nil {
		err = fmt.Errorf("Unable to parse arch from package string: %s", err)
	}
	return
}

// Pops the release from the end of an RPM string with architecture
// already popped off.
func popRelease(noArch string) (popped string, release string, err error) {
	popped, release, err = popDelim(noArch, "-")
	if err != nil {
		err = fmt.Errorf("Unable to parse release from package string: %s", err)
	}
	return
}

// Pops the version from the end of an RPM string with arch and release already
// popped off
func popVersion(noArchRelease string) (popped string, version string, err error) {
	popped, version, err = popDelim(noArchRelease, "-")
	if err != nil {
		err = fmt.Errorf("Unable to parse version from package string: %s", err)
	}
	return
}

// Parse the Epoch and package name from what's left in the RPM string
func parseEpochName(inStr string) (epoch string, name string) {
	split := strings.Split(inStr, ":")
	if len(split) == 1 {
		epoch = "0"
		name = split[0]
		return
	}
	epoch = split[0]
	name = split[1]
	return
}

// Generic function for popping the last group from a string given a variable
// delimiter.
func popDelim(inStr string, delim string) (outStr string, remain string, err error) {
	split := strings.Split(inStr, delim)
	if len(split) == 1 {
		err = fmt.Errorf("Splitting %s at %s produced one result", inStr, delim)
		return
	}
	outStr = split[len(split)-1]
	remain = strings.Split(inStr, delim+outStr)[0]
	return
}
