package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

type Package struct {
	Path      string
	BuildInfo PkgBuild
	SkipDebug bool
	config    *Config
}

func (config *Config) NewPackage(name string) (*Package, error) {
	packagePath := filepath.Join(config.AurPath, name)

	return &Package{
		Path:   packagePath,
		config: config,
	}, nil
}

func (pkg *Package) Cleanup() error {
	logs, err := filepath.Glob(filepath.Join(pkg.Path, "*.log"))
	if err != nil {
		return err
	}

	for _, log := range logs {
		if err = os.Remove(log); err != nil {
			return err
		}
	}

	return nil
}

func (pkg *Package) ArchivesExistInRepo() bool {
	archive := pkg.FileNames()[0]
	if _, err := os.Stat(filepath.Join(pkg.config.RepoPath, archive)); err != nil {
		return false
	}

	return true
}

func (pkg *Package) OutOfDate() bool {
	if pkg.BuildInfo.Version.Current != pkg.BuildInfo.Version.Latest {
		return true
	}

	return !pkg.ArchivesExistInRepo()
}

func (pkg *Package) MoveArchivesToRepo() error {
	for i, filename := range pkg.FileNames() {
		if err := moveFile(filepath.Join(pkg.Path, filename), filepath.Join(pkg.config.RepoPath, filename)); err != nil {
			if i == 1 {
				pkg.SkipDebug = true
			} else {
				return err
			}
		}
	}

	return nil
}

func (pkg *Package) UpdateRepo() error {
	if _, err := os.Stat(pkg.config.RepoPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			os.MkdirAll(pkg.config.RepoPath, 0660)
		} else {
			return err
		}
	}

	fmt.Println("\x1b[1;34m  ->\x1b[0m \x1b[1;37mMoving archives\x1b[0m")
	if err := pkg.MoveArchivesToRepo(); err != nil {
		return err
	}

	fmt.Println("\x1b[1;34m  ->\x1b[0m \x1b[1;37mSigning archives\x1b[0m")
	if err := pkg.SignArchives(); err != nil {
		return err
	}

	return pkg.AddToRepo()
}

func moveFile(source, target string) error {
	if err := os.Rename(source, target); err == nil {
		return nil
	}

	fileInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	bar := progressbar.DefaultBytes(fileInfo.Size(), fileInfo.Name())

	if _, err = io.Copy(io.MultiWriter(targetFile, bar), sourceFile); err != nil {
		return err
	}

	if err := targetFile.Sync(); err != nil {
		return err
	}

	if err := os.Chmod(target, fileInfo.Mode()); err != nil {
		return err
	}

	return os.Remove(source)
}
