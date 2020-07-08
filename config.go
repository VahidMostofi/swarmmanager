package swarmmanager

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// config ...
type config struct {
	ResultsDirectoryPath string `yaml:"resultsPath"`
}

var c *config

func (c *config) check() {
	c.ResultsDirectoryPath = strings.Trim(c.ResultsDirectoryPath, " ")
	fi, err := os.Stat(c.ResultsDirectoryPath)
	if err != nil {
		panic(err)
	}
	if !fi.Mode().IsDir() {
		panic(fmt.Errorf("the path is not a directory: %s", c.ResultsDirectoryPath))
	}
	if c.ResultsDirectoryPath[len(c.ResultsDirectoryPath)-1] != '/' {
		c.ResultsDirectoryPath += "/"
	}
	fmt.Println("Config: ResultPath is", c.ResultsDirectoryPath)
}

// GetConfig ...
func GetConfig() *config {
	if c == nil {
		c = &config{}
		b, e := ioutil.ReadFile("config.yml")
		if e != nil {
			panic(e)
		}
		e = yaml.Unmarshal(b, c)
		if e != nil {
			panic(e)
		}
		c.check()
	}
	return c
}

func init() {
	GetConfig()
}
