package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	App struct {
		Name string `json:"name"`
	} `json:"app"`
	Web struct {
		Prefork bool `json:"prefork"`
		Port    int  `json:"port"`
	} `json:"web"`
	Log struct {
		Level int `json:"level"`
	} `json:"log"`
	PostgreSQL struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
		SSLMode  string `json:"sslmode"`
		TimeZone string `json:"timezone"`
	} `json:"postgresql"`
	MongoDB struct {
		URI string `json:"uri"`
	} `json:"mongodb"`
	Redis struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
	RabbitMQ struct {
		URI string `json:"uri"`
	} `json:"rabbitmq"`
	JWTSecret string `json:"jwt_secret"`
}

var AppConfig Config

func LoadConfig() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	fmt.Println("Environment: ", env)

	// Log the current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get current working directory: %v", err)
	}
	log.Printf("Current working directory: %s", wd)

	// Open the config file
	fileName := "config." + env + ".json"
	configFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Cannot open config file: %v", err)
	}
	defer configFile.Close()

	// Decode the config file
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.Fatalf("Cannot decode config JSON: %v", err)
	}

	log.Println("Configuration loaded successfully")
}
