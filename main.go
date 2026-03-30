package main

import (
	"fmt"
	"os"
)

func mustReadDir(name string) []os.DirEntry {
	entries, err := os.ReadDir(name)
	if err != nil {
		panic(err.Error())
	}
	return entries
}

func main() {
	config := NewConfig()
	for _, entry := range mustReadDir(config.AurPath) {
		if !entry.IsDir() {
			continue
		}

		fmt.Print("\x1b[1;34m==>\x1b[0m \x1b[1;37mChecking ", entry.Name(), "\x1b[0m")

		pkg, err := config.NewPackage(entry.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}

		if pkg.CheckAURUpdated() {
			if err := pkg.UpdateAURRepo(); err != nil {
				panic(err.Error())
			}
		}

		if err := pkg.ParsePKGBUILD(); err != nil {
			panic(err.Error())
		}

		if !pkg.OutOfDate() {
			fmt.Println(" - Up to date")
			continue
		}

		fmt.Println("\n\x1b[1;34m==>\x1b[0m \x1b[1;37mBuilding\x1b[0m")
		if err := pkg.Build(); err != nil {
			fmt.Println(" - Build fialed", err)
			continue
		}

		fmt.Println("\x1b[1;34m==>\x1b[0m \x1b[1;37mUpdating repo\x1b[0m")
		if err := pkg.UpdateRepo(); err != nil {
			panic(err.Error())
		}
		fmt.Println(" - done")

		fmt.Print("\x1b[1;34m==>\x1b[0m \x1b[1;37mCleanup\x1b[0m")
		if err := pkg.Cleanup(); err != nil {
			panic(err.Error())
		}
		fmt.Println(" - done")
	}
}
