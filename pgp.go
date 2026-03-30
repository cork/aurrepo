package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func (pkg *Package) SignArchives() error {
	for _, filename := range pkg.FileNames() {
		if _, err := os.Stat(filepath.Join(pkg.config.RepoPath, filename)); err != nil {
			continue
		}

		cmd := exec.Command("gpg", "--detach-sign", "--no-armor", filename)
		cmd.Dir = pkg.config.RepoPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
