package main

type Configuration struct {
	Databases []Database `yaml:"databases"`
}

type Database struct {
	ConnectUrl  string       `yaml:"urlConnect"`
	Name        string       `yaml:"name"`
	Collections []Collection `yaml:"collections"`
}

type Collection struct {
	Name    string   `yaml:"name"`
	Indexes []string `yaml:"indexes"`
}
