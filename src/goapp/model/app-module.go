package model

type ApplicationModule struct {
	ApplicationId         string `json:"applicationId"`
	ApplicationName       string `json:"applicationName"`
	ApplicationModuleId   string `json:"applicationModuleId"`
	ApplicationModuleName string `json:"applicationModuleName"`
	Callbackurl           string `json:"callbackurl"`
	RequireRemarks        bool   `json:"requireRemarks"`
	ApproveText           string `json:"approveText"`
	RejectText            string `json:"rejectText"`
	ReassignCallbackUrl   string `json:"reassignCallbackUrl"`
	ExportUrl             string `json:"exportUrl"`
	AllowReassign         bool   `json:"allowReassign"`
	RequireAuthentication bool   `json:"requireAuthentication"`
}
