package model

import "html/template"

type ApprovalRequestEmailData struct {
	Subject     string
	Body        template.HTML
	ApproveText string
	RejectText  string
	ApproveUrl  string
	RejectUrl   string
}
