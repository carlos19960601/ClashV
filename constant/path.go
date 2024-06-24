package constant

import (
	"os"
	P "path"
)

const Name = "clash"

var Path = func() *path {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir, _ = os.Getwd()
	}

	homeDir = P.Join(homeDir, ".config", Name)
	if _, err = os.Stat(homeDir); err != nil {

	}
	return &path{
		homeDir:    homeDir,
		configFile: "config.yaml",
	}
}()

type path struct {
	homeDir    string
	configFile string
}

func SetHomeDir(root string) {
	Path.homeDir = root
}

func SetConfig(file string) {
	Path.configFile = file
}

func (p *path) HomeDir() string {
	return p.homeDir
}

func (p *path) Config() string {
	return p.configFile
}
