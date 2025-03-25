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

func (ecm *envConfigManager) GetCallbackRetryFreq() string {
	return os.Getenv("CALLBACK_RETRY_FREQ")
}

func (ecm *envConfigManager) GetClientID() string {
	return os.Getenv("CLIENT_ID")
}

func (ecm *envConfigManager) GetClientSecret() string {
	return os.Getenv("CLIENT_SECRET")
}

func (ecm *envConfigManager) GetCommunityPortalAppId() string {
	return os.Getenv("COMMUNITY_PORTAL_APP_ID")
}

func (ecm *envConfigManager) GetCommunityPortalDomain() string {
	return os.Getenv("COMMUNITY_PORTAL_DOMAIN")
}

func (ecm *envConfigManager) GetDatabaseConnectionString() string {
	return os.Getenv("APPROVALSYSTEMDB_CONNECTION_STRING")
}

func (ecm *envConfigManager) GetEmailClientID() string {
	return os.Getenv("EMAIL_CLIENT_ID")
}

func (ecm *envConfigManager) GetEmailClientSecret() string {
	return os.Getenv("EMAIL_CLIENT_SECRET")
}

func (ecm *envConfigManager) GetCTO() string {
	return os.Getenv("CTO")
}

func (ecm *envConfigManager) GetEmailTenantID() string {
	return os.Getenv("EMAIL_TENANT_ID")
}

func (ecm *envConfigManager) GetEmailUserID() string {
	return os.Getenv("EMAIL_USER_ID")
}

func (ecm *envConfigManager) GetEnterpriseOwners() []string {
	enterpriseOwners := os.Getenv("ENTERPRISE_OWNERS")
	if enterpriseOwners == "" {
		return nil
	}
	ownersArray := strings.Split(enterpriseOwners, ",")
	return ownersArray
}

func (ecm *envConfigManager) GetHomeURL() string {
	return os.Getenv("HOME_URL")
}

func (ecm *envConfigManager) GetIPDRAppId() string {
	return os.Getenv("IPDR_APP_ID")
}

func (ecm *envConfigManager) GetIPDRModuleId() string {
	return os.Getenv("IPDR_MODULE_ID")
}

func (ecm *envConfigManager) GetIsEmailEnabled() bool {
	return os.Getenv("EMAIL_ENABLED") == "true"
}

func (ecm *envConfigManager) GetLinkFooters() string {
	return os.Getenv("LINK_FOOTERS")
}

func (ecm *envConfigManager) GetOrganizationName() string {
	return os.Getenv("ORGANIZATION_NAME")
}

func (ecm *envConfigManager) GetPort() string {
	return os.Getenv("PORT")
}

func (ecm *envConfigManager) GetScope() string {
	return os.Getenv("SCOPE")
}

func (ecm *envConfigManager) GetSessionKey() string {
	return os.Getenv("SESSION_KEY")
}

func (ecm *envConfigManager) GetTenantID() string {
	return os.Getenv("TENANT_ID")
}
