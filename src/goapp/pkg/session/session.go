package session

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"main/pkg/sql"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	auth "main/pkg/authentication"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	Store *sessions.FilesystemStore
)

func InitializeSession() {
	Store = sessions.NewFilesystemStore(os.TempDir(), []byte("secret"))
	Store.MaxLength(math.MaxInt64)
	gob.Register(map[string]interface{}{})
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	var url string
	url = fmt.Sprintf("%v", r.URL)
	if strings.HasPrefix(url, "/response/") {
		authReq := true
		params := mux.Vars(r)
		var appModuleGuid string
		appModuleGuid = params["appModuleGuid"]

		authReq = IsAuthRequired(appModuleGuid)
		if authReq == false {
			return true
		}
	}

	// Check session if there is saved user profile
	url = fmt.Sprintf("/loginredirect?redirect=%v", r.URL)
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		c := http.Cookie{
			Name:   "auth-session",
			MaxAge: -1}
		http.SetCookie(w, &c)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return false
	}
	// fmt.Println(session)
	if _, ok := session.Values["profile"]; !ok {

		// Asks user to login if there is no saved user profile
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return false

	} else {
		// If there is a user profile saved
		authenticator, err := auth.NewAuthenticator()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}

		// Retrieve current token data
		refreshToken := fmt.Sprintf("%s", session.Values["refresh_token"])
		expiry, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s", session.Values["expiry"]))
		today := time.Now().UTC()

		if today.After(expiry) {
			// Attempt to refresh token if access token is already expired
			ts := authenticator.Config.TokenSource(context.Background(), &oauth2.Token{RefreshToken: refreshToken, Expiry: expiry})
			newToken, err := ts.Token()
			if err != nil {
				// fmt.Printf("ERROR REFRESHING TOKEN\n")
				if rErr, ok := err.(*oauth2.RetrieveError); ok {
					details := new(ErrorDetails)
					if err := json.Unmarshal(rErr.Body, details); err != nil {
						panic(err)
					}

					fmt.Println(details.Error, details.ErrorDescription)
				}

				// Log out the user if the attempt to refresh the token failed
				http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)

			} else if newToken != nil {
				// fmt.Printf("TOKEN REFRESHED\n")

				// Save new token data
				session.Values["refresh_token"] = newToken.RefreshToken
				session.Values["expiry"] = newToken.Expiry.UTC().Format("2006-01-02 15:04:05")
				session.Values["access_token"] = newToken.AccessToken
				err = session.Save(r, w)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return false
				}
			}
		}
		return true
	}
}

func IsGuidAuthenticated(w http.ResponseWriter, r *http.Request) {
	// Check header if authenticated
	_, err := auth.VerifyAccessToken(r)
	// RETURN ERROR
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetState(w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	if _, ok := session.Values["state"]; ok {

		return fmt.Sprintf("%s", session.Values["state"]), nil

	}
	return "", nil
}

func connectSql() (db *sql.DB) {
	db, err := sql.Init(sql.ConnectionParam{ConnectionString: os.Getenv("APPROVALSYSTEMDB_CONNECTION_STRING")})
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
	}
	return
}

func IsAuthRequired(appModuleGuid string) bool {
	db := connectSql()
	defer db.Close()
	queryParams := map[string]interface{}{
		"ApplicationModuleId": appModuleGuid,
	}
	isAuthRequired := true
	res, err := db.ExecuteStoredProcedureWithResult("PR_ApplicationModules_IsAuthRequired", queryParams)
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
	} else {
		if res[0] != nil {
			isAuthRequired = res[0]["RequireAuthentication"].(bool)
		}
	}
	return isAuthRequired
}

type ErrorDetails struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
