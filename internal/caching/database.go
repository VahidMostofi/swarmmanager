package caching

import (
	"crypto/md5"
	"fmt"
	"sort"

	"github.com/VahidMostofi/swarmmanager/internal/autoconfigure"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/VahidMostofi/swarmmanager"
)

// Database interface for caching mechanism
type Database interface {
	Store(string, map[string]swarm.ServiceSpecs, autoconfigure.Information) error
	Retrieve(string, map[string]swarm.ServiceSpecs) autoconfigure.Information
}

// MongoDatabase ...
type MongoDatabase struct {
	Client *mongo.Client
	Col    *mongo.Collection
}

// GetNewMongoDatabaseConnector ...
func GetNewMongoDatabaseConnector() (*Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(swarmmanager.GetConfig().MongoDBURL))
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongodb: %w", err)
	}
	collectionName := "configs"
	col, err := client.Database("configcache").Collection(collectionName)
	if err != nil {
		return nil, fmt.Errorf("error connecting to collection %s: %w", collectionName, err)
	}
	md := &MongoDatabase{
		Client: client,
		Col:    col,
	}

	return md, nil
}

// Store ...
func (md *MongoDatabase) Store(string, map[string]swarm.ServiceSpecs, autoconfigure.Information) error {
	return nil
}

// Retrieve ...
func (md *MongoDatabase) Retrieve(workload string, configs map[string]swarm.ServiceSpecs) autoconfigure.Information {
	fmt.Println(md.hash(workload, configs))
	return nil
}

func (md *MongoDatabase) hash(workload string, configs map[string]swarm.ServiceSpecs) string {
	bytes := make([]byte, 0)
	bytes = append(bytes, []byte("v1")...) // TODO get from config
	bytes = append(bytes, []byte(workload)...)

	var keys []string
	var tempConfigs = make(map[string]swarm.ServiceSpecs)
	for key, value := range configs {
		tempConfigs[value.ImageName] = value
		keys = append(keys, value.ImageName)
	}

	sort.Strings(keys)
	for _, key := range keys {
		fmt.Println(key)
		bytes = append(bytes, tempConfigs[key].GetBytes()...)
	}

	return fmt.Sprintf("%x", md5.Sum(bytes))
}
