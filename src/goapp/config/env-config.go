package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type envConfigManager struct {
	*Config
}

func NewEnvConfigManager() *envConfigManager {
	// Set environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	return &envConfigManager{}
}

func (ecm *envConfigManager) GetDatabaseConnectionString() string {
	return os.Getenv("APPROVALSYSTEMDB_CONNECTION_STRING")
}

func (ecm *envConfigManager) GetEnterpriseOwners() []string {
	enterpriseOwners := os.Getenv("ENTERPRISE_OWNERS")
	if enterpriseOwners == "" {
		return nil
	}
	ownersArray := strings.Split(enterpriseOwners, ",")
	return ownersArray
}
