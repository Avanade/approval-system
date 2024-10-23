package config

type Key string

type Config struct {
	DatabaseConnectionString string
}

type ConfigManager interface {
	GetDatabaseConnectionString() string
	GetEnterpriseOwners() []string
	GetEmailTenantID() string
	GetEmailClientID() string
	GetEmailClientSecret() string
	GetEmailUserID() string
	GetHomeURL() string
	GetIsEmailEnabled() bool
	GetTenantID() string
	GetClientID() string
	GetClientSecret() string
	GetLinkFooters() string
	GetOrganizationName() string
	GetCommunityPortalAppId() string
	GetCallbackRetryFreq() string
	GetPort() string
}
