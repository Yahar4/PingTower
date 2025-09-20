package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBname   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

// ConnectDB connects to a DataBase with specified
// .yaml configuration using dsn
func ConnectPostgresDB(logger *zap.SugaredLogger) (*sqlx.DB, error) {
	config, err := loadConfig("configs/database.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Database.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.DBname,
		config.Database.Port,
		config.Database.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	logger.Info("Successfully connected to database")
	return db, nil
}

// loadConfig made to unmarshal the .yaml file
// with specified path
func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, err
}
