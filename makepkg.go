package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (pkg *Package) FileNames() []string {
	cmd := exec.Command("makepkg", "--packagelist")
	cmd.Dir = pkg.Path
	cmd.Stderr = os.Stderr

	data, err := cmd.Output()
	if err != nil {
		panic(err.Error())
	}

	archive := strings.Replace(string(data[0:len(data)-1]), pkg.Path+"/", "", 1)

	archive = strings.Replace(archive, fmt.Sprintf("%s-%s", pkg.BuildInfo.Name, pkg.BuildInfo.Version.Current), fmt.Sprintf("%s-%s", pkg.BuildInfo.Name, pkg.BuildInfo.Version.Latest), 1)

	return []string{
		archive,
		strings.Replace(archive, pkg.BuildInfo.Name, pkg.BuildInfo.Name+"-debug", 1),
	}
}
