package item

type ReassignItemCallback struct {
	Id                  string `json:"Id"`
	ApproverEmail       string `json:"ApproverEmail"`
	Username            string `json:"Username"`
	ApplicationId       string `json:"ApplicationId"`
	ApplicationModuleId string `json:"ApplicationModuleId"`
	ApproveText         string `json:"ApproveText"`
	RejectText          string `json:"RejectText"`
}
