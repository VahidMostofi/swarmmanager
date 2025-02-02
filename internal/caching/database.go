package caching

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/history"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"gopkg.in/yaml.v3"
)

// Database interface for caching mechanism
type Database interface {
	Store(string, map[string]swarm.ServiceSpecs, history.Information) (string, error)
	Retrieve(string, map[string]swarm.ServiceSpecs) (history.Information, error)
	GetNotFoundError() error
}

// DropboxDatabase ...
type DropboxDatabase struct {
	Path string
}

// GetNewDropboxDatabase ...
func GetNewDropboxDatabase(path string) (Database, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("error checking Dropbox path %s: %w", path, err)
	}
	if !fi.Mode().IsDir() {
		return nil, fmt.Errorf("the Dropbox path is not a directory: %s", path)
	}
	return &DropboxDatabase{
		Path: path,
	}, nil
}

// GetNotFoundError ...
func (md *DropboxDatabase) GetNotFoundError() error {
	return fmt.Errorf("not found")
}

// Store ...
func (md *DropboxDatabase) Store(workload string, specs map[string]swarm.ServiceSpecs, info history.Information) (string, error) {
	hash := md.hash(workload, specs)
	info.HashCode = hash
	log.Println("DropboxCache: Storing: hash for this configuration and workload:", hash)
	b, err := yaml.Marshal(info)
	if err != nil {
		return "", fmt.Errorf("DropboxCache: error while converting information to yaml: %w", err)
	}
	os.Mkdir(md.Path+"/"+hash, 1777)
	err = ioutil.WriteFile(md.Path+"/"+hash+"/info.yml", b, 0777)
	if err != nil {
		return "", fmt.Errorf("DropboxCache: error while storing information file on %s: %w", md.Path+"/"+hash+"info.yml", err)
	}
	return hash, nil
}

// Retrieve ...
func (md *DropboxDatabase) Retrieve(workload string, specs map[string]swarm.ServiceSpecs) (history.Information, error) {
	hash := md.hash(workload, specs)
	log.Println("DropboxCache: Retrieving: hash for this configuration and workload:", hash)
	h := history.Information{}
	if _, err := os.Stat(md.Path + "/" + hash + "/info.yml"); os.IsNotExist(err) {
		log.Println("DropboxCache: result for this configuration/workload not found in cache")
		return h, md.GetNotFoundError()
	}
	b, err := ioutil.ReadFile(md.Path + "/" + hash + "/info.yml")
	if err != nil {
		return h, fmt.Errorf("DropboxCache: the done file exists at %s but can't read info.yml: %w", md.Path+"/"+hash, err)
	}
	yaml.Unmarshal(b, &h)
	return h, nil
}

func (md *DropboxDatabase) hash(workload string, specs map[string]swarm.ServiceSpecs) string {
	code := configs.GetConfig().Version + "_" + configs.GetConfig().AppName + "_" + workload + "_" + strconv.Itoa(configs.GetConfig().Test.Duration) + "_"
	cpus := ""
	var keys []string
	for _, srv := range configs.GetConfig().TestBed.ServicesToConfigure {
		keys = append(keys, srv)
	}
	sort.Strings(keys)
	for _, srv := range keys {
		count := specs[srv].CPULimits * float64(specs[srv].ReplicaCount)
		countStr := strconv.FormatFloat(count, 'f', 1, 64)
		cpus += countStr + "_"
	}
	code += cpus
	return code
}

// func (md *DropboxDatabase) hash(workload string, specs map[string]swarm.ServiceSpecs) string {
// 	bytes := make([]byte, 0)
// 	bytes = append(bytes, []byte(configs.GetConfig().Version)...)
// 	bytes = append(bytes, []byte(configs.GetConfig().AppName)...)
// 	bytes = append(bytes, []byte(workload)...)
// 	bytes = append(bytes, []byte(strconv.Itoa(configs.GetConfig().Test.Duration))...)
// 	str := ""
// 	for key, value := range configs.GetConfig().LoadGenerator.Args {
// 		str += value + key
// 	}
// 	bytes = append(bytes, []byte(str)...)

// 	var keys []string
// 	var tempConfigs = make(map[string]swarm.ServiceSpecs)
// 	for _, value := range specs {
// 		tempConfigs[value.Name] = value
// 		keys = append(keys, value.Name)
// 	}

// 	sort.Strings(keys)
// 	for _, key := range keys {
// 		// fmt.Println("hash with", tempConfigs[key])
// 		bytes = append(bytes, tempConfigs[key].GetBytes()...)
// 	}
// 	// fmt.Println("hash with", configs.GetConfig().Version, configs.GetConfig().AppName, workload)
// 	return fmt.Sprintf("%x", md5.Sum(bytes))
// }
