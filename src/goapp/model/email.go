package model

type TypEmailMessage struct {
	To      string
	Subject string
	Body    string
}

type TypEmailData struct {
	Subject     string
	Body        string
	ApproveText string
	RejectText  string
	ApproveUrl  string
	RejectUrl   string
}
