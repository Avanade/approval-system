package model

import "html/template"

type TypEmailMessage struct {
	To      string
	Subject string
	Body    string
}

type TypEmailData struct {
	Subject     string
	Body        template.HTML
	ApproveText string
	RejectText  string
	ApproveUrl  string
	RejectUrl   string
}
