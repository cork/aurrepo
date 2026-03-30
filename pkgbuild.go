package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type PkgBuild struct {
	Name    string
	Version PkgVersion
	Release uint8
	Epoch   uint8
	Arch    []string
}

type PkgVersion struct {
	HasFunction bool
	Current     string
	Latest      string
}

var _arch = regexp.MustCompile(`(?m)^\s*arch=\(([^)]+)\)`)
var _epoch = regexp.MustCompile(`(?m)^\s*epoch=(\d+)$`)
var _pkgname = regexp.MustCompile(`(?m)^\s*pkgname=(.+)$`)
var _pkgrel = regexp.MustCompile(`(?m)^\s*pkgrel=(\d+)$`)
var _pkgver = regexp.MustCompile(`(?m)^\s*pkgver=(.+)$`)
var _pkgverfunc = regexp.MustCompile(`(?m)^\s*pkgver\(\)\s*\{$`)

func (pkg *Package) ParsePKGBUILD() error {
	file, err := os.Open(filepath.Join(pkg.Path, "PKGBUILD"))
	if err != nil {
		return err
	}
	defer file.Close()

	pkgbuild, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if match := _pkgname.FindSubmatch(pkgbuild); len(match) > 1 {
		pkg.BuildInfo.Name = string(match[1])
	}

	if match := _pkgver.FindSubmatch(pkgbuild); len(match) > 1 {
		pkg.BuildInfo.Version.Current = string(match[1])
		pkg.BuildInfo.Version.Latest = pkg.BuildInfo.Version.Current
	}

	pkg.BuildInfo.Version.HasFunction = _pkgverfunc.Match(pkgbuild)

	if match := _arch.FindSubmatch(pkgbuild); len(match) > 1 {
		pkg.BuildInfo.Arch = strings.Split(strings.ReplaceAll(strings.ReplaceAll(string(match[1]), "'", ""), "\n", " "), " ")
	}

	if match := _pkgrel.FindSubmatch(pkgbuild); len(match) > 1 {
		if _, err := fmt.Sscanf(string(match[1]), "%d", &pkg.BuildInfo.Release); err != nil {
			return err
		}
	}

	if match := _epoch.FindSubmatch(pkgbuild); len(match) > 1 {
		if _, err := fmt.Sscanf(string(match[1]), "%d", &pkg.BuildInfo.Epoch); err != nil {
			return err
		}
	}

	if pkg.BuildInfo.Version.HasFunction {
		pkg.BuildInfo.Version.Latest = pkg.LatestVersion()
	}

	return nil
}

func (pkg *Package) LatestVersion() string {
	cmd := exec.Command("makepkg", "--noprepare", "-od")
	cmd.Dir = pkg.Path
	if err := cmd.Run(); err != nil {
		panic(err.Error())
	}

	cmd = exec.Command("bash", "-c", "source ../PKGBUILD; pkgver")
	cmd.Dir = filepath.Join(pkg.Path, "src")
	cmd.Env = append(os.Environ(), "srcdir="+filepath.Join(pkg.Path, "src"))
	cmd.Stderr = os.Stderr
	data, err := cmd.Output()
	if err != nil {
		panic(err.Error())
	}

	os.RemoveAll(filepath.Join(pkg.Path, "src"))

	return strings.Trim(string(data), "\t\n ")
}
