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
