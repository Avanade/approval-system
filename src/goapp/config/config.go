package config

type Key string

type Config struct {
	DatabaseConnectionString string
}

type ConfigManager interface {
	GetCallbackRetryFreq() string
	GetClientID() string
	GetClientSecret() string
	GetCTO() string
	GetCommunityPortalAppId() string
	GetCommunityPortalDomain() string
	GetDatabaseConnectionString() string
	GetEmailClientID() string
	GetEmailClientSecret() string
	GetEmailTenantID() string
	GetEmailUserID() string
	GetEnterpriseOwners() []string
	GetHomeURL() string
	GetIPDRAppId() string
	GetIPDRModuleId() string
	GetIsEmailEnabled() bool
	GetLinkFooters() string
	GetOrganizationName() string
	GetPort() string
	GetScope() string
	GetSessionKey() string
	GetTenantID() string
}
