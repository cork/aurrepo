package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func (pkg *Package) AddToRepo() error {
	args := []string{
		"--new",
		"--remove",
		"--include-sigs",
		"--sign",
		filepath.Join(pkg.config.RepoPath, pkg.config.DBName),
	}

	for i, archive := range pkg.FileNames() {
		if i == 1 && pkg.SkipDebug {
			continue
		}

		args = append(args, archive)
	}

	cmd := exec.Command("repo-add", args...)
	cmd.Dir = pkg.config.RepoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
