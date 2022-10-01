package lib

import (
	"Tavern-Backend/models"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Configuration Struct For Proper Security
type Configuration struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		Username string `yaml:"username"`
		Database string `yaml:"database"`
	} `yaml:"database"`
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Cors struct {
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
		v.SetConfigName(".local")
	} else {
		v.SetConfigName(".prod")
	}
	v.SetConfigType("yaml")

	// set the configuration file path
	// it should be in ./env
	v.AddConfigPath("./TavernProfile-env")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s:%s\n", v.GetString("server_host"), v.GetString("server_port"))

	var config Configuration
	var m map[string]interface{}
	err = v.Unmarshal(&m)
	if err != nil {
		panic(err)
	}

	err = mapstructure.Decode(m, &config)
	if err != nil {
		for _, i := range v.AllKeys() {
			fmt.Printf("%s\n", i)
		}
		fmt.Print("\n")
		panic(err)
	}

	fmt.Printf("%v\n", config)

	return config
}

func (config Configuration) GetDatabaseConnectionString(local bool) string {
	if local {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Database)
	} else {
		return fmt.Sprintf("%s:%s@tcp(172.17.0.1:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.Database.Username, config.Database.Password, config.Database.Port, config.Database.Database)
	}
}

func (config Configuration) GetEmailConfig() models.AuthEmailConfiglette {
	return models.AuthEmailConfiglette{
		Host:     config.Email.Host,
		Port:     config.Email.Port,
		Username: config.Email.Username,
		Password: config.Email.Password,
	}
}
