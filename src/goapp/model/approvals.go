package model

type TypRequestProcess struct {
	ApplicationId       string `json:"applicationId"`
	ApplicationModuleId string `json:"applicationModuleId"`
	ItemId              string `json:"itemId"`
	ApproverEmail       string `json:"approverEmail"`
	Remarks             string `json:"remarks"`
	IsApproved          string `json:"isApproved"`
	Username            string `json:"username"`
}
