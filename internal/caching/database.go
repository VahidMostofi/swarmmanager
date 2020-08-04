package caching

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/VahidMostofi/swarmmanager"
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
func (md *DropboxDatabase) Store(workload string, configs map[string]swarm.ServiceSpecs, info history.Information) (string, error) {
	hash := md.hash(workload, configs)
	info.HashCode = hash
	log.Println("DropboxCache: hash for this configuration and workload:", hash)
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
func (md *DropboxDatabase) Retrieve(workload string, configs map[string]swarm.ServiceSpecs) (history.Information, error) {
	hash := md.hash(workload, configs)
	log.Println("DropboxCache: hash for this configuration and workload:", hash)
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

func (md *DropboxDatabase) hash(workload string, configs map[string]swarm.ServiceSpecs) string {
	bytes := make([]byte, 0)
	bytes = append(bytes, []byte(swarmmanager.GetConfig().Version)...)
	bytes = append(bytes, []byte(swarmmanager.GetConfig().SystemName)...)
	bytes = append(bytes, []byte(workload)...)

	var keys []string
	var tempConfigs = make(map[string]swarm.ServiceSpecs)
	for _, value := range configs {
		tempConfigs[value.ImageName] = value
		keys = append(keys, value.ImageName)
	}

	sort.Strings(keys)
	for _, key := range keys {
		bytes = append(bytes, tempConfigs[key].GetBytes()...)
	}

	return fmt.Sprintf("%x", md5.Sum(bytes))
}
