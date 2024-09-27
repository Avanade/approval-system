package model

type Item struct {
	Id                  string   `json:"id"`
	Application         string   `json:"application"`
	ApproverRemarks     string   `json:"approverRemarks"`
	Body                string   `json:"body"`
	Created             string   `json:"created"`
	DateResponded       string   `json:"dateResponded"`
	DateSent            string   `json:"dateSent"`
	IsApproved          bool     `json:"isApproved"`
	Module              string   `json:"module"`
	Subject             string   `json:"subject"`
	ApproveText         string   `json:"approveText"`
	CallbackUrl         string   `json:"callbackUrl"`
	ReassignCallbackUrl string   `json:"reassignCallbackUrl"`
	RejectText          string   `json:"rejectText"`
	ApproveUrl          string   `json:"approveUrl"`
	RejectUrl           string   `json:"rejectUrl"`
	AllowReassign       bool     `json:"allowReassign"`
	AllowReassignUrl    string   `json:"allowReassignUrl"`
	RespondedBy         string   `json:"respondedBy"`
	Approvers           []string `json:"approvers"`
	RequestedBy         string   `json:"requestedBy"`
}

type ItemOptions struct {
	ItemType     int64
	ItemStatus   int64
	Offset       int
	Filter       int
	Search       string
	RequestType  string
	Organization string
	User         string
}

type Response struct {
	Data  []Item `json:"data"`
	Total int    `json:"total"`
}

type ItemType int8

const (
	RequestItem ItemType = iota
	ApprovalItem
	AllType
)

type ItemStatus int8

const (
	Pending ItemStatus = iota
	Approved
	Rejected
	Closed // Disapproved, Approved
	All    // Disapproved, Approved, Pending
)

type ItemInsertRequest struct {
	ApplicationId       string   `json:"applicationId"`
	ApplicationModuleId string   `json:"applicationModuleId"`
	RequesterEmail      string   `json:"requesterEmail"`
	Emails              []string `json:"emails"`
	Subject             string   `json:"subject"`
	Body                string   `json:"body"`
}

type ApprovalRequestApprover struct {
	ItemId        string
	ApproverEmail string
}

type ProcessResponseRequest struct {
	ApplicationId       string `json:"applicationId"`
	ApplicationModuleId string `json:"applicationModuleId"`
	ItemId              string `json:"itemId"`
	ApproverEmail       string `json:"approverEmail"`
	Remarks             string `json:"remarks"`
	IsApproved          string `json:"isApproved"`
	Username            string `json:"username"`
}

type ResponseCallback struct {
	ItemId       string `json:"itemId"`
	IsApproved   bool   `json:"isApproved"`
	Remarks      string `json:"remarks"`
	ResponseDate string `json:"responseDate"`
	RespondedBy  string `json:"respondedBy"`
}
