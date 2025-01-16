package configs

import (
	"fmt"
	"os"
)

type Settings struct {
	DatabaseDSN string
	SecretJWT   []byte
	ServiceHost string
	ServicePort string
	ServiceProtocol string
}

func NewSettings() (*Settings, error) {

	// Postgres
	postgres_user := os.Getenv("POSTGRES_USER")
	postgres_password := os.Getenv("POSTGRES_PASSWORD")
	postgres_db := os.Getenv("POSTGRES_DB")
	postgres_host := os.Getenv("POSTGRES_HOST")
	postgres_port := os.Getenv("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", postgres_host, postgres_port,
		postgres_user, postgres_db, postgres_password)
	// Service
	SERVICE_HOST := os.Getenv("SERVICE_HOST")
	SERVICE_PORT := os.Getenv("SERVICE_PORT")
	SERVICE_PROTOCOL := os.Getenv("SERVICE_PROTOCOL")

	return &Settings{
		DatabaseDSN: dsn,
		SecretJWT:   []byte(os.Getenv("SECRET_JWT")),
		ServiceHost: SERVICE_HOST,
		ServicePort: SERVICE_PORT,
		ServiceProtocol: SERVICE_PROTOCOL,
	}, nil
}

func (s Settings) String() string {
	result := fmt.Sprintf("Settings:\n\tServerHost:%s\n\tServerPort:%s\n\tServiceProtocol:%s\n\tDatabaseDSN:%s\n",
			              s.ServiceHost, s.ServicePort, s.ServiceProtocol, s.DatabaseDSN)
	return result
}

var SettingsGlobal, _ = NewSettings()
