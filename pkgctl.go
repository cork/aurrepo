package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func (pkg *Package) Build() error {
	args := []string{"build", "-c"}
	if list, ok := pkg.config.Install[pkg.BuildInfo.Name]; ok {
		for _, pack := range list {
			match, err := filepath.Glob(filepath.Join(pkg.config.RepoPath, pack))
			if err != nil {
				return err
			}

			if len(match) < 1 {
				return fmt.Errorf("package %s not found", pack)
			}

			args = append(args, "-I", match[0])
		}
	}

	cmd := exec.Command("pkgctl", args...)
	cmd.Dir = pkg.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
