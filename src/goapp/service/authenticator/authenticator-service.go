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
	Session     *session.Session
}

func NewAuthenticatorService(conf config.ConfigManager, session *session.Session) AuthenticatorService {
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
	s, err := a.Session.Get(r, session)
	if err != nil {
		return err
	}

	s.Options.MaxAge = -1
	err = s.Save(r, *w)
	if err != nil {
		return err
	}

	return nil
}

func (a *authenticatorService) GetAuthCodeURL(state string) string {
	return a.OAuthConfig.AuthCodeURL(state)
}

func (a *authenticatorService) GetAuthenticatedUser(r *http.Request) (*model.AzureUser, error) {
	s, err := a.Session.Get(r, "auth-session")
	if err != nil {
		return nil, err
	}

	var profile map[string]interface{}
	u := s.Values["profile"]
	profile, ok := u.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error getting user data")
	}

	return &model.AzureUser{
		Name:  profile["name"].(string),
		Email: profile["preferred_username"].(string),
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
	s, err := a.Session.Get(r, session)
	if err != nil {
		return "", err
	}

	if _, ok := s.Values[key]; !ok {
		return "", fmt.Errorf("no key found in session")
	}

	return fmt.Sprintf("%s", s.Values[key]), nil
}

func (a *authenticatorService) IsAuthenticated(r *http.Request) (bool, bool, error) {
	s, err := a.Session.Get(r, "auth-session")
	if err != nil {
		return false, false, err
	}

	if _, ok := s.Values["profile"]; !ok {
		return false, false, fmt.Errorf("no profile found in session")
	}

	today := time.Now().UTC()
	expiry, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s", s.Values["expiry"]))
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
	s, err := a.Session.Get(r, "auth-session")
	if err != nil {
		return err
	}

	refreshToken := fmt.Sprintf("%s", s.Values["refresh_token"])
	expiry, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s", s.Values["expiry"]))
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
		s.Values["refresh_token"] = newToken.RefreshToken
		s.Values["expiry"] = newToken.Expiry.UTC().Format("2006-01-02 15:04:05")
		s.Values["access_token"] = newToken.AccessToken
		err = s.Save(r, *w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *authenticatorService) SaveOnSession(w *http.ResponseWriter, r *http.Request, session string, p ...interface{}) error {
	s, err := a.Session.Get(r, session)
	if err != nil {
		return err
	}

	for _, v := range p {
		switch v := v.(type) {
		case model.SessionMapValue:
			s.Values[v.Key] = v.Value
		case model.SessionStringValue:
			s.Values[v.Key] = v.Value
		}
	}

	s.Options.MaxAge = 43200
	err = s.Save(r, *w)
	if err != nil {
		return err
	}

	return nil
}
