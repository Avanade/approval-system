package model

type SessionStringValue struct {
	Key   string
	Value string
}

type SessionMapValue struct {
	Key   string
	Value map[string]interface{}
}

type AzureUser struct {
	Name            string `json:"name"`
	Email           string `json:"preferred_username"`
	IsLegalApprover bool
	IsAuditor       bool
}
