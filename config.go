package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

type Config struct {
	Arch     string
	AurPath  string
	RepoPath string
	DBName   string
	Install  map[string][]string
}

func NewConfig() *Config {
	config := Config{}

	pflag.StringVar(&config.Arch, "arch", "x86_64", "Arch to update repo with")
	pflag.StringVarP(&config.AurPath, "aur", "a", "~/.cache/aurrepo/aur", "Path to where the aur packages are located")
	pflag.StringVarP(&config.RepoPath, "repo", "r", "~/.cache/aurrepo/repo", "Path to repo directory")
	pflag.StringVarP(&config.DBName, "dbname", "d", "aurrepo.db.tar.zst", "Repo db name")
	help := pflag.Bool("help", false, "display this help and exit")

	pflag.Parse()

	if *help {
		pflag.Usage()
		os.Exit(1)
	}

	config.AurPath = strings.ReplaceAll(config.AurPath, "~", "${HOME}")
	config.RepoPath = strings.ReplaceAll(config.RepoPath, "~", "${HOME}")

	config.AurPath = os.ExpandEnv(config.AurPath)
	config.RepoPath = os.ExpandEnv(config.RepoPath)
	config.RepoPath = filepath.Join(config.RepoPath, config.Arch)

	if _, err := os.Stat(filepath.Join(config.AurPath, "install.json")); err == nil {
		file, err := os.Open(filepath.Join(config.AurPath, "install.json"))
		if err != nil {
			panic(err.Error())
		}
		defer file.Close()

		if err := json.NewDecoder(file).Decode(&config.Install); err != nil {
			panic(err.Error())
		}
	}

	return &config
}
