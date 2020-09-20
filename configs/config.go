package configs

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var c *Configurations

// Configurations ...
type Configurations struct {
	Version        string
	AppName        string
	Jaeger         JaegerConfigurations
	Cache          CacheConfigurations
	Results        ResultsConfigurations
	TestBed        TestBedConfigurations
	Host           HostConfigurations
	Log            LogConfigurations
	LoadGenerator  LoadGeneratorConfigurations
	Test           TestConfigurations
	UsageCollector UsageCollectorConfigurations
}

// TestConfigurations ...
type TestConfigurations struct {
	Duration                   int
	WaitAfterLoadGeneratorDone int
}

// LoadGeneratorConfigurations ...
type LoadGeneratorConfigurations struct {
	Type    string
	Details map[string]string
	Args    map[string]string
}

// LogConfigurations ...
type LogConfigurations struct {
	Directory string
}

// HostConfigurations ...
type HostConfigurations struct {
	AvailableCPUCount float64
	Host              string
}

// TestBedConfigurations ...
type TestBedConfigurations struct {
	DockerComposeFile   string
	ServiceCount        int
	StackName           string
	ServicesToConfigure []string
}

// ResultsConfigurations ...
type ResultsConfigurations struct {
	Path string
}

// JaegerConfigurations ...
type JaegerConfigurations struct {
	Host            string
	RootService     string
	DetailsFilePath string
	StorePath       string
}

// CacheConfigurations ...
type CacheConfigurations struct {
	Type    string
	Details map[string]string
}

// UsageCollectorConfigurations ...
type UsageCollectorConfigurations struct {
	Type    string
	Details map[string]interface{}
}

// GetConfig ...
func GetConfig() *Configurations {
	if c == nil {
		panic(fmt.Errorf("the configuration should have been initialized before this"))
	}
	return c
}

// Initialize ...
func Initialize() {
	c = &Configurations{}
	err := viper.Unmarshal(c)
	if err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile(c.Log.Directory+"/"+time.Now().Local().Format(time.RFC3339)+".log", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	c.AppName = viper.GetString("appName")

	if strings.Contains(GetConfig().Results.Path, "$WORKLOAD") {
		GetConfig().Results.Path = strings.Replace(GetConfig().Results.Path, "$WORKLOAD", viper.GetString("workloadStr"), 1)
		log.Println("Updating result path to", GetConfig().Results.Path)
	}

	if strings.Contains(GetConfig().Results.Path, "$SYSTEM_NAME") {
		GetConfig().Results.Path = strings.Replace(GetConfig().Results.Path, "$SYSTEM_NAME", GetConfig().AppName, 1)
		log.Println("Updating result path to", GetConfig().Results.Path)
	}

}
