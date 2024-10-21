package authenticator

type UserToken struct {
	IdToken      string
	AccessToken  string
	Profile      map[string]interface{}
	RefreshToken string
	Expiry       string
}
