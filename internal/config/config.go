package config

import (
	"errors"
	"log" // Pour logger les informations ou erreurs de chargement de config

	"github.com/spf13/viper" // La bibliothèque pour la gestion de configuration
)

type Config struct {
	Server struct {
		Port    int    `mapstructure:"port"`
		BaseURL string `mapstructure:"base_url"`
	} `mapstructure:"server"`
	Database struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"database"`
	Analytics struct {
		BufferSize  int `mapstructure:"buffer_size"`
		WorkerCount int `mapstructure:"worker_count"`
	} `mapstructure:"analytics"`
	Monitor struct {
		IntervalMinutes int `mapstructure:"interval_minutes"`
	} `mapstructure:"monitor"`
}

func LoadConfig() (*Config, error) {

	var cfg Config

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.SetConfigType("yaml")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.base_url", "http://localhost:8080")
	viper.SetDefault("database.name", "default_db")
	viper.SetDefault("analytics.buffer_size", 100)
	viper.SetDefault("monitor.interval_minutes", 5)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("Fichier de configuration 'config.yaml' non trouvé, utilisation des valeurs par défaut ou des variables d'environnement.")
		}
	} else {
		log.Println("Fichier de configuration 'config.yaml' chargé avec succès.")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Impossible de désérialiser la configuration: %v", err)
	}

	log.Printf("Configuration loaded: Server Port=%d, DB Name=%s, Analytics Buffer=%d, Monitor Interval=%dmin",
		cfg.Server.Port, cfg.Database.Name, cfg.Analytics.BufferSize, cfg.Monitor.IntervalMinutes)

	return &cfg, nil
}
