package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

const configFile = ".gatorconfig.json"

type Config struct {
	DBURL       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
	Filepointer *os.File
	Path        string
}

func ReadConfig() *Config {
	homedir, _ := os.UserHomeDir()
	path := filepath.Join(homedir, configFile)
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Fatal("unable to open config")
	}

	decoder := json.NewDecoder(io.Reader(jsonFile))
	gatorConfig := &Config{}
	if err = decoder.Decode(gatorConfig); err != nil {
		log.Fatal("cannot decode JSON into config struct")
	}
	gatorConfig.Filepointer = jsonFile
	gatorConfig.Path = path

	return gatorConfig
}

func (cfg *Config) SetUser(username string) {
	cfg.CurrentUser = username

	data, err := json.Marshal(cfg)
	if err != nil {
		log.Fatal("unable to marshal json")
	}

	os.WriteFile(cfg.Path, data, 0644)

}
