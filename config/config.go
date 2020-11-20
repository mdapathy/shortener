package config

import (
	"fmt"
	"github.com/tkanos/gonfig"
	"log"
	"net/url"
)

type Configuration struct {
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	Database string `json:"Database"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

//Creates database configuration based on config.json
func New(configFilename string) Configuration {
	config := Configuration{}
	err := gonfig.GetConf(configFilename, &config)
	if err != nil {
		log.Fatalf("Cannot instantiate server:\t%s", err)
	}
	return config
}

//Generates a database connection url based on given configuration
func (c *Configuration) ConnectionUrl() string {
	dbUrl := &url.URL{
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%s", c.Host, c.Port),
		User:   url.UserPassword(c.Username, c.Password),
		Path:   c.Database,
	}

	return dbUrl.String()
}
