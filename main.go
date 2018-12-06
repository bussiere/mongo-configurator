package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
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
			var connectionString string
			var dataBaseNameProvided string

			envVarConnStr, exist := isThatEnv(dataBaseUrl)
			if !exist {
				connectionString = dataBaseUrl
			} else {
				connectionString = getEnv(envVarConnStr)
			}

			if len(connectionString) == 0 {
				log.Printf("dataBaseUrl was not provided")
				return
			}

			envVarDbName, exist := isThatEnv(dataBaseName)
			if !exist {
				dataBaseNameProvided = dataBaseName
			} else {
				dataBaseNameProvided = getEnv(envVarDbName)
			}

			if len(dataBaseNameProvided) == 0 {
				log.Printf("dataBaseName was not provided")
				return
			}

			conStr := strings.Replace(connectionString, "+", "%2B", -1);
			sessionMongo, err := mgo.Dial(conStr)
			if err != nil {
				log.Printf("Can't set mongo connected to %s\n", conStr)
				return
			}
			defer sessionMongo.Close()

			for _, collectionConfig := range collectionArray {
				var collectionNameProvided string

				envCollectionName, exist := isThatEnv(collectionConfig.Name)

				if !exist {
					collectionNameProvided = collectionConfig.Name
				} else {
					collectionNameProvided = getEnv(envCollectionName)
				}

				if len(collectionNameProvided) == 0 {
					log.Printf("CollectionName was not provided")
					return
				}

				collection := sessionMongo.DB(dataBaseNameProvided).C(collectionNameProvided)
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
						log.Printf("Can't create index for database=%s, collection=%s, index=%s", dataBaseNameProvided, collectionNameProvided, indexName)
					} else {
						log.Printf("Created index for database=%s, collection=%s, index=%s", dataBaseNameProvided, collectionNameProvided, indexName)
					}
				}

				log.Printf("Provide to dataBase=%s, collection=%s indexes", dataBaseNameProvided, collectionNameProvided)
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

func isThatEnv(str string) (string, bool) {
	re := regexp.MustCompile("\\$\\{(.*?)\\}")
	match := re.FindStringSubmatch(str)
	if len(match) > 1 {
		return match[1], true
	}
	return "", false
}

func getEnv(str string) string {
	env := os.Getenv(str)
	if len(env) == 0 {
		return ""
	}
	log.Printf("Parse env var %s", env)
	return env
}
