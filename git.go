package main

import (
	"os"
	"os/exec"
)

func (pkg *Package) CheckAURUpdated() bool {
	cmd := exec.Command("git", "fetch")
	cmd.Dir = pkg.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err.Error())
	}

	cmd = exec.Command("git", "log", "-1", "--pretty=format:%H")
	cmd.Dir = pkg.Path
	cmd.Stderr = os.Stderr

	data, err := cmd.Output()
	if err != nil {
		panic(err.Error())
	}
	currentReference := string(data)

	cmd = exec.Command("git", "log", "origin/master", "-1", "--pretty=format:%H")
	cmd.Dir = pkg.Path
	cmd.Stderr = os.Stderr

	data, err = cmd.Output()
	if err != nil {
		panic(err.Error())
	}
	originReference := string(data)

	return currentReference != originReference
}

func (pkg *Package) UpdateAURRepo() error {
	cmd := exec.Command("git", "reset", "--hard", "@{u}")
	cmd.Dir = pkg.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}
