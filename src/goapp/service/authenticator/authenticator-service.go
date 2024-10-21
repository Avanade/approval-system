package authenticator

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"main/pkg/envvar"
	"net/http"
	"os"
	"strings"

	oidc "github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"golang.org/x/oauth2"
)

type authenticatorService struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticatorService() AuthenticatorService {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://login.microsoftonline.com/"+os.Getenv("TENANT_ID")+"/v2.0")
	if err != nil {
		log.Printf("failed to get provider: %v", err)
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  envvar.GetEnvVar("HOME_URL", "http://localhost:8080") + "/login/azure/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile"},
	}

	return &authenticatorService{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
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

func (a *authenticatorService) GetConfig() *oauth2.Config {
	return &a.Config
}
