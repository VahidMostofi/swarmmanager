package swarmmanager

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// config ...
type config struct {
	ResultsDirectoryPath string `yaml:"resultsPath"`
	LogDirectory         string `yaml:"logDirectory"`
	JaegerStorePath      string `yaml:"jegerStorePath"`
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
	log.Println("Config: ResultPath is", c.ResultsDirectoryPath)
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
		logFile, err := os.OpenFile(c.LogDirectory+"/"+time.Now().Local().Format(time.RFC3339)+".log", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
		c.check()
	}
	return c
}

func init() {
	GetConfig()
}
