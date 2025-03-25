package model

type IPDRequest struct {
	RequestId         int64    `json:"requestId"`
	RequestorName     string   `json:"requestorName"`
	RequestorEmail    string   `json:"requestorEmail"`
	IPTitle           string   `json:"title"`
	IPType            string   `json:"type"`
	IPDescription     string   `json:"description"`
	Reason            string   `json:"reason"`
	InvolvementId     []string `json:"involvementId"`
	Involvement       []string `json:"involvement"`
	IsApproved        bool     `json:"isApproved"`
	ApprovalRequestId string   `json:"approvalRequestId"`
	ApproverRemarks   string   `json:"approverRemarks"`
	ResponseDate      string   `json:"responseDate"`
	RespondedBy       string   `json:"respondedBy"`
	Created           string   `json:"created"`
}
