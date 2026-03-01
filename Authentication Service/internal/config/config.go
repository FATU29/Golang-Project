package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUrl string

	DbHost     string
	DbUser     string
	DbPassword string
	DbPort     string
	DbName     string
	DbType     string

	Port      string
	Host      string
	JwtSecret string

	BrevoApiKey       string
	SenderName        string
	SenderMail        string
	BrevoEndpointMail string

	RedisHost     string
	RedisPassword string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DbUrl:             os.Getenv("DATABASE_URL"),
		DbHost:            os.Getenv("DB_HOST"),
		DbUser:            os.Getenv("DB_USER"),
		DbPassword:        os.Getenv("DB_PASSWORD"),
		DbPort:            os.Getenv("DB_PORT"),
		DbName:            os.Getenv("DB_NAME"),
		DbType:            os.Getenv("DB_TYPE"),
		Port:              os.Getenv("PORT"),
		JwtSecret:         os.Getenv("JWT_SECRET"),
		Host:              os.Getenv("HOST"),
		BrevoEndpointMail: os.Getenv("BREVO_SENDMAIL_ENDPOINT"),
		BrevoApiKey:       os.Getenv("BREVO_API_KEY"),
		SenderName:        os.Getenv("SENDER_NAME"),
		SenderMail:        os.Getenv("SENDER_MAIL"),
		RedisHost:         os.Getenv("REDIS_HOST"),
		RedisPassword:     os.Getenv("REDIS_PASSWORD"),
		GoogleClientID:    os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL: os.Getenv("GOOGLE_REDIRECT_URL"),
	}
}

func (c *Config) GetDbUrl() string {
	return c.DbUrl
}

func (c *Config) GetDbHost() string {
	return c.DbHost
}

func (c *Config) GetDbUser() string {
	return c.DbUser
}

func (c *Config) GetDbPassword() string {
	return c.DbPassword
}

func (c *Config) GetDbName() string {
	return c.DbName
}

func (c *Config) GetDbType() string {
	return c.DbType
}
func (c *Config) GetPort() string {
	return c.Port
}

func (c *Config) GetHost() string {
	return c.Host
}

func (c *Config) GetJwtSecret() string {
	return c.JwtSecret
}
