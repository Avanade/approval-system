package model

type LegalConsultation struct {
	ItemId    string `json:"itemId"`
	Email     string `json:"email"`
	Created   string `json:"created"`
	CreatedBy string `json:"createdBy"`
}

type ConsultLegalRequest struct {
	ItemId              string `json:"itemId"`
	ApplicationId       string `json:"applicationId"`
	ApplicationModuleId string `json:"applicationModuleId"`
}

type Approver struct {
	ApprovalTypeId int    `json:"approvalTypeId"`
	ApproverEmail  string `json:"approverEmail"`
	ApproverName   string `json:"approverName"`
}
