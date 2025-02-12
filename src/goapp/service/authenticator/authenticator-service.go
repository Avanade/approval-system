package authenticator

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"main/config"
	"main/infrastructure/session"
	"main/model"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"golang.org/x/oauth2"
)

type authenticatorService struct {
	Provider    *oidc.Provider
	OAuthConfig oauth2.Config
	Ctx         context.Context
	Config      config.ConfigManager
	Session     session.ConnectSession
}

func NewAuthenticatorService(conf config.ConfigManager, session session.ConnectSession) AuthenticatorService {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://login.microsoftonline.com/"+conf.GetTenantID()+"/v2.0")
	if err != nil {
		log.Printf("failed to get provider: %v", err)
	}

	oauth2Conf := oauth2.Config{
		ClientID:     conf.GetClientID(),
		ClientSecret: conf.GetClientSecret(),
		RedirectURL:  conf.GetHomeURL() + "/login/azure/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile"},
	}

	return &authenticatorService{
		Provider:    provider,
		OAuthConfig: oauth2Conf,
		Ctx:         ctx,
		Config:      conf,
		Session:     session,
	}
}

func (a *authenticatorService) AccessTokenIsValid(r *http.Request) bool {
	tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	keySet, err := jwk.Fetch(r.Context(), "https://login.microsoftonline.com/common/discovery/v2.0/keys")
	if err != nil {
		log.Println(err.Error())
		return false
	}

	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwa.RS256.String() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %v not found", kid)
		}

		publickey := &rsa.PublicKey{}
		err = keys.Raw(publickey)
		if err != nil {
			return nil, fmt.Errorf("could not parse pubkey")
		}

		return publickey, nil
	})

	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

func (a *authenticatorService) ClearFromSession(w *http.ResponseWriter, r *http.Request, session string) error {
	s := a.Session.ConnectSession(w, r, session)
	err := s.Destroy()
	if err != nil {
		return err
	}

	return nil
}

func (a *authenticatorService) GenerateToken() (string, error) {
	urlPath := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", a.Config.GetTenantID())
	client := &http.Client{
		Timeout: time.Second * 90,
	}

	data := url.Values{}
	data.Set("client_id", a.Config.GetClientID())
	data.Set("scope", a.Config.GetScope())
	data.Set("client_secret", a.Config.GetClientSecret())
	data.Set("grant_type", "client_credentials")
	ecodedData := data.Encode()

	req, err := http.NewRequest("POST", urlPath, strings.NewReader(ecodedData))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var token struct {
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		ExtExpiresIn int    `json:"ext_expires_in"`
		AccessToken  string `json:"access_token"`
	}
	err = json.NewDecoder(response.Body).Decode(&token)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func (a *authenticatorService) GetAuthCodeURL(state string) string {
	return a.OAuthConfig.AuthCodeURL(state)
}

func (a *authenticatorService) GetAuthenticatedUser(r *http.Request) (*model.AzureUser, error) {
	//Get Azure user profile
	var profile map[string]interface{}
	s := a.Session.ConnectSession(nil, r, "auth-session")
	u, err := s.Get("profile")
	if err != nil {
		return nil, err
	}

	profile, ok := u.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error getting user data")
	}

	// Check if user is a legal approver
	isLegalApprover := false
	b, err := s.Get("isLegalApprover")
	if err == nil {
		isLegalApprover, _ = b.(bool)
	}

	return &model.AzureUser{
		Name:            profile["name"].(string),
		Email:           profile["preferred_username"].(string),
		IsLegalApprover: isLegalApprover,
	}, nil
}

func (a *authenticatorService) GetLogoutURL() (string, error) {
	logoutUrl, err := url.Parse("https://login.microsoftonline.com/" + a.Config.GetTenantID() + "/oauth2/logout?client_id=" + a.Config.GetClientID() + "&post_logout_redirect_uri=" + a.Config.GetHomeURL())
	if err != nil {
		return "", err
	}
	return logoutUrl.String(), nil
}

func (a *authenticatorService) GetStringValue(r *http.Request, session string, key string) (string, error) {
	s := a.Session.ConnectSession(nil, r, session)
	v, err := s.Get(key)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", v), nil
}

func (a *authenticatorService) IsAuthenticated(r *http.Request) (bool, bool, error) {
	s := a.Session.ConnectSession(nil, r, "auth-session")

	_, err := s.Get("profile")
	if err != nil {
		return false, false, err
	}

	e, err := s.Get("expiry")
	if err != nil {
		return false, false, err
	}

	today := time.Now().UTC()
	expiry, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s", e))
	if err != nil {
		return false, false, err
	}

	return true, today.After(expiry), nil
}

func (a *authenticatorService) ProcessToken(code string) (*UserToken, error) {
	token, err := a.OAuthConfig.Exchange(context.TODO(), code)
	if err != nil {
		log.Printf("no token found: %v", err)
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.Config.GetClientID(),
	}

	idToken, err := a.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

	if err != nil {
		return nil, fmt.Errorf("failed to verify id token: %v", err)
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return nil, err
	}

	return &UserToken{
		IdToken:      rawIDToken,
		AccessToken:  token.AccessToken,
		Profile:      profile,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.UTC().Format("2006-01-02 15:04:05"),
	}, nil
}

func (a *authenticatorService) RefreshToken(w *http.ResponseWriter, r *http.Request) error {
	s := a.Session.ConnectSession(w, r, "auth-session")

	t, err := s.Get("refresh_token")
	if err != nil {
		return err
	}

	e, err := s.Get("expiry")
	if err != nil {
		return err
	}

	refreshToken := fmt.Sprintf("%s", t)
	expiry, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s", e))
	if err != nil {
		return err
	}

	ts := a.OAuthConfig.TokenSource(context.Background(), &oauth2.Token{RefreshToken: refreshToken, Expiry: expiry})
	newToken, err := ts.Token()
	if err != nil {
		if rErr, ok := err.(*oauth2.RetrieveError); ok {
			details := new(ErrorDetails)
			if err := json.Unmarshal(rErr.Body, details); err != nil {
				fmt.Println(err)
				return err
			}

			fmt.Println(details.Error, details.ErrorDescription)
			return fmt.Errorf("error refreshing token: %v", details.ErrorDescription)
		}
	}

	if newToken != nil {
		s.Set("refresh_token", newToken.RefreshToken)
		s.Set("expiry", newToken.Expiry.UTC().Format("2006-01-02 15:04:05"))
		s.Set("access_token", newToken.AccessToken)
		err = s.Save(43200)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *authenticatorService) SaveOnSession(w *http.ResponseWriter, r *http.Request, session string, p map[string]interface{}) error {
	s := a.Session.ConnectSession(w, r, session)

	for k, val := range p {
		s.Set(k, val)
	}

	err := s.Save(43200)
	if err != nil {
		return err
	}

	return nil
}
