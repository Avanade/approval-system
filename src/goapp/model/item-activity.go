package model

type ItemActivity struct {
	Id          int64  `json:"id"`
	CreatedBy   string `json:"createdBy"`
	Created     string `json:"created"`
	Content     string `json:"content"`
	ItemId      string `json:"itemId"`
	AppId       string `json:"appId"`
	AppModuleId string `json:"appModuleId"`
}

type InvolvedUsers struct {
	Requestor   string
	Approvers   []string
	Consultants []string
}

type Activity struct {
	Action  string       `json:"action"`
	Created string       `json:"created"`
	Actor   string       `json:"actor"`
	Comment ItemActivity `json:"details"`
}
