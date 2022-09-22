package lib

import (
	"Tavern-Backend/models"

	"github.com/spf13/viper"
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
	// use viper to load the configuration file
	v := viper.New()

	// set the configuration file name
	if local {
		v.SetConfigName("local")
	} else {
		v.SetConfigName("prod")
	}
	v.SetConfigType("yaml")

	// set the configuration file path
	// it should be in ./env
	v.AddConfigPath("./env")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config Configuration
	err = v.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func (config Configuration) GetDatabaseConnectionString() string {
	return config.Database.Username + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + config.Database.Port + ")/" + config.Database.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func (config Configuration) GetEmailConfig() models.AuthEmailConfiglette {
	return models.AuthEmailConfiglette{
		Host:     config.Email.Host,
		Port:     config.Email.Port,
		Username: config.Email.Username,
		Password: config.Email.Password,
	}
}
