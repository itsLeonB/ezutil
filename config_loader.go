package ezutil

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GenericConfig interface {
	Prefix() string
}

// Config represents the complete application configuration loaded from environment variables.
// It aggregates all configuration sections including app settings, authentication, database parameters,
// and the initialized GORM DB instance.
type Config struct {
	App     *App
	Auth    *Auth
	SQLDB   *SQLDB
	GORM    *gorm.DB
	Generic GenericConfig
}

// LoadConfig reads environment variables into the default Config, loads sub-configuration
// for App, Auth, and SQLDB, establishes a GORM connection, and returns the fully initialized Config.
func LoadConfig(defaults Config) *Config {
	return LoadConfigWithDB(defaults, true)
}

// LoadConfigWithDB reads environment variables into the default Config with optional database connection.
// If connectDB is false, the GORM field will be nil, allowing for testing without database dependency.
func LoadConfigWithDB(defaults Config, connectDB bool) *Config {
	appDefaults := App{}
	if defaults.App != nil {
		appDefaults = *defaults.App
	}

	authDefaults := Auth{}
	if defaults.Auth != nil {
		authDefaults = *defaults.Auth
	}

	config := &Config{
		App:  loadAppConfig(appDefaults),
		Auth: loadAuthConfig(authDefaults),
	}

	if connectDB {
		sqlDBConfig := loadSQLDBConfig()
		config.SQLDB = sqlDBConfig
		config.GORM = sqlDBConfig.openGormConnection()
		if defaults.Generic != nil {
			config.Generic = loadGenericConfig(defaults.Generic.Prefix(), defaults.Generic)
		}
	} else {
		// For testing, load SQLDB config without validation/connection
		config.SQLDB = loadSQLDBConfigOptional()
	}

	return config
}

// LoadConfigWithoutDB loads configuration without establishing database connection.
// This is useful for testing or when database connection is not needed.
func LoadConfigWithoutDB(defaults Config) *Config {
	return LoadConfigWithDB(defaults, false)
}

// App holds application-level settings such as environment name, server port, request timeout,
// allowed client URLs, and timezone. These values are populated from environment variables.
type App struct {
	Env        string
	Port       string
	Timeout    time.Duration
	ClientUrls []string
	Timezone   string
}

func loadAppConfig(defaults App) *App {
	var loadedConfig App

	err := envconfig.Process("APP", &loadedConfig)
	if err != nil {
		log.Fatalf("error loading app config: %s", err.Error())
	}

	if loadedConfig.Env == "" {
		loadedConfig.Env = defaults.Env
	}
	if loadedConfig.Port == "" {
		loadedConfig.Port = defaults.Port
	}
	// Validate port number (only if port is provided)
	if loadedConfig.Port != "" {
		if port, err := strconv.Atoi(loadedConfig.Port); err != nil || port < 1 || port > 65535 {
			log.Fatalf("invalid port number: %s", loadedConfig.Port)
		}
	}
	if loadedConfig.Timeout <= 0 {
		if loadedConfig.Timeout < 0 {
			log.Println("timeout cannot be negative, using default value...")
		}
		loadedConfig.Timeout = defaults.Timeout
	}
	if len(loadedConfig.ClientUrls) == 0 {
		loadedConfig.ClientUrls = defaults.ClientUrls
	}
	if loadedConfig.Timezone == "" {
		loadedConfig.Timezone = defaults.Timezone
	}
	// Validate timezone: TODO
	// if _, err := time.LoadLocation(loadedConfig.Timezone); err != nil {
	// 	log.Fatalf("invalid timezone: %s", loadedConfig.Timezone)
	// }

	return &loadedConfig
}

// Auth holds authentication configuration including JWT secret key, token and cookie durations,
// issuer identifier, and authentication service URL. Values are sourced from environment variables.
type Auth struct {
	SecretKey      string
	TokenDuration  time.Duration
	CookieDuration time.Duration
	Issuer         string
	URL            string
}

func loadAuthConfig(defaults Auth) *Auth {
	var loadedConfig Auth

	err := envconfig.Process("AUTH", &loadedConfig)
	if err != nil {
		log.Fatalf("error loading auth config: %s", err.Error())
	}

	if loadedConfig.SecretKey == "" {
		loadedConfig.SecretKey = defaults.SecretKey
	}
	if loadedConfig.TokenDuration <= 0 {
		if loadedConfig.TokenDuration < 0 {
			log.Println("token duration cannot be negative, using default value...")
		}
		loadedConfig.TokenDuration = defaults.TokenDuration
	}
	if loadedConfig.CookieDuration <= 0 {
		if loadedConfig.CookieDuration < 0 {
			log.Println("cookie duration cannot be negative, using default value...")
		}
		loadedConfig.CookieDuration = defaults.CookieDuration
	}
	if loadedConfig.Issuer == "" {
		loadedConfig.Issuer = defaults.Issuer
	}
	if loadedConfig.URL == "" {
		loadedConfig.URL = defaults.URL
	}

	return &loadedConfig
}

// SQLDB holds SQL database connection parameters loaded from environment variables,
// including host, user credentials, database name, port, and driver type for GORM.
type SQLDB struct {
	Host     string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true"`
	Name     string `required:"true"`
	Port     string `required:"true"`
	Driver   string `required:"true"`
}

func loadSQLDBConfig() *SQLDB {
	var loadedConfig SQLDB

	err := envconfig.Process("SQLDB", &loadedConfig)
	if err != nil {
		log.Fatalf("error loading SQLDB config: %s", err.Error())
	}

	// Only validate port if it's provided (for testing scenarios)
	if loadedConfig.Port != "" {
		if port, err := strconv.Atoi(loadedConfig.Port); err != nil || port < 1 || port > 65535 {
			log.Fatalf("invalid database port number: %s", loadedConfig.Port)
		}
	}

	return &loadedConfig
}

// loadSQLDBConfigOptional loads SQLDB config without failing if required fields are missing.
// This is useful for testing scenarios where database config is not needed.
func loadSQLDBConfigOptional() *SQLDB {
	var loadedConfig SQLDB

	// Process without failing on missing required fields
	_ = envconfig.Process("SQLDB", &loadedConfig)

	// Only validate port if it's provided
	if loadedConfig.Port != "" {
		if port, err := strconv.Atoi(loadedConfig.Port); err != nil || port < 1 || port > 65535 {
			log.Printf("warning: invalid database port number: %s", loadedConfig.Port)
		}
	}

	return &loadedConfig
}

func (sqldb *SQLDB) openGormConnection() *gorm.DB {
	db, err := gorm.Open(sqldb.getGormDialector(), &gorm.Config{})
	if err != nil {
		log.Fatalf("error opening GORM connection: %s", err.Error())
	}

	return db
}

func (sqldb *SQLDB) getGormDialector() gorm.Dialector {
	switch sqldb.Driver {
	case "mysql":
		return mysql.Open(sqldb.getDSN())
	case "postgres":
		return postgres.Open(sqldb.getDSN())
	default:
		log.Fatalf("unsupported SQLDB driver: %s", sqldb.Driver)
		return nil
	}
}

func (sqldb *SQLDB) getDSN() string {
	switch sqldb.Driver {
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			sqldb.User,
			sqldb.Password,
			sqldb.Host,
			sqldb.Port,
			sqldb.Name,
		)
	case "postgres":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s",
			sqldb.Host,
			sqldb.User,
			sqldb.Password,
			sqldb.Name,
			sqldb.Port,
		)
	default:
		log.Fatalf("unsupported SQLDB driver: %s", sqldb.Driver)
		return ""
	}
}

func loadGenericConfig[T GenericConfig](prefix string, defaults T) T {
	if err := envconfig.Process(prefix, defaults); err != nil {
		log.Fatalf("error loading generic config: %s", err.Error())
	}

	return defaults
}
