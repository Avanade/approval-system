package authenticator

type UserToken struct {
	IdToken      string
	AccessToken  string
	Profile      map[string]interface{}
	RefreshToken string
	Expiry       string
}

type ErrorDetails struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
