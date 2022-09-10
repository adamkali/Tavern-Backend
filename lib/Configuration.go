package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Configuration Struct For Proper Security
type Configuration struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		Username string `yaml:"username"`
		Database string `yaml:"database"`
	} `yaml:"database"`
	ServerPort string `yaml:"server_port"`
	ServerHost string `yaml:"server_host"`
	Cors       struct {
		AllowedOrigins []string `yaml:"origins"`
		AllowedMethods []string `yaml:"methods"`
		AllowedHeaders []string `yaml:"headers"`
		Credentials    bool     `yaml:"credentials"`
	} `yaml:"cors"`
	Email struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"email"`
}

func LoadConfiguration(local bool) Configuration {
	// if local then load .local.config
	// else load the .prod.config

	// get the operating system
	// if windows set the locPath to .\\env\\.local.config
	//and prodPath to  .\\env\\.prod.config
	// else set the locPath to ./env/.local.config
	// and prodPath to ./env/.prod.config
	var locPath string
	var prodPath string
	if runtime.GOOS == "windows" {
		locPath, _ = filepath.Abs(".\\env\\.local.config")
		prodPath, _ = filepath.Abs(".\\env\\.prod.config")
	} else {
		locPath, _ = filepath.Abs("./env/.local.config")
		prodPath, _ = filepath.Abs("./env/.prod.config")
	}

	var config Configuration

	if local {
		//Print local
		fmt.Print("Loading Local Configuration\n")
		// get .local.config from the /env folder
		// check the operating system and load the correct file.
		yamlfile, err := os.Open(locPath)
		if err != nil {
			panic(err)
		}
		yamlbytes, err := ioutil.ReadAll(yamlfile)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlbytes, &config)
		if err != nil {
			println(err.Error())
			panic(err)
		}
		fmt.Print(config)
	} else {
		yamlfile, err := os.Open(prodPath)
		if err != nil {
			panic(err)
		}
		yamlbytes, err := ioutil.ReadAll(yamlfile)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlbytes, &config)
		if err != nil {
			panic(err)
		}
	}
	return config
}

func (config Configuration) GetDatabaseConnectionString() string {
	return config.Database.Username + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + config.Database.Port + ")/" + config.Database.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
}
