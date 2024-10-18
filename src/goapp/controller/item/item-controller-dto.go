package item

import "main/model"

type ReassignItemCallback struct {
	Id                  string `json:"Id"`
	ApproverEmail       string `json:"ApproverEmail"`
	Username            string `json:"Username"`
	ApplicationId       string `json:"ApplicationId"`
	ApplicationModuleId string `json:"ApplicationModuleId"`
	ApproveText         string `json:"ApproveText"`
	RejectText          string `json:"RejectText"`
}

type RespondePageData struct {
	ApplicationId       string
	ApplicationModuleId string
	ItemId              string
	ApproverEmail       string
	IsApproved          string
	Data                model.Item
	RequireRemarks      bool
	Response            string
	ApproveText         string
	RejectText          string
}
