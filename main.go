package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("Config file must be first argument")
	}
	args := os.Args[1:]
	configFile := args[0]
	if len(configFile) == 0 {
		log.Fatalf("Config file must be first argument")
	}

	configuration := parseFile(configFile)
	var waitGroup sync.WaitGroup
	for _, dataBase := range configuration.Databases {
		waitGroup.Add(1)
		go func(dataBaseUrl string, dataBaseName string, collectionArray []Collection) {
			defer waitGroup.Done()
			var err error
			conStr := strings.Replace(dataBaseUrl, "+", "%2B", -1);
			sessionMongo, err := mgo.Dial(conStr)
			if err != nil {
				log.Printf("Can't set mongo connected to %s\n", conStr)
				return
			}
			defer sessionMongo.Close()

			for _, collectionConfig := range collectionArray {
				collection := sessionMongo.DB(dataBaseName).C(collectionConfig.Name)
				for _, indexName := range collectionConfig.Indexes {
					index := mgo.Index{
						Key:        []string{indexName},
						Unique:     true,
						DropDups:   true,
						Background: true,
						Sparse:     true,
					}

					err := collection.EnsureIndex(index)
					if err != nil {
						log.Printf("Can't create index for database=%s, collection=%s, index=%s", dataBaseName, collectionConfig.Name, indexName)
					} else {
						log.Printf("Created index for database=%s, collection=%s, index=%s", dataBaseName, collectionConfig.Name, indexName)
					}
				}

				log.Printf("Provide to dataBase=%s, collection=%s indexes", dataBaseName, collectionConfig.Name)
			}
		}(dataBase.ConnectUrl, dataBase.Name, dataBase.Collections)
	}

	waitGroup.Wait()
	log.Printf("Configration success applied")
}

func parseFile(file string) Configuration {
	var configuration Configuration
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Can't read file %s", file)
	}
	err = yaml.Unmarshal(yamlFile, &configuration)
	if err != nil {
		log.Fatalf("Can't parse file %s, with error %s", file, err)
	}

	return configuration
}
