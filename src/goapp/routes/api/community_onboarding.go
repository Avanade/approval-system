package routes

import (
	"encoding/json"
	db "main/pkg/ghmgmtdb"
	session "main/pkg/session"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetCommunityOnBoardingInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 0, 64)

	// Get email address of the user
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	switch r.Method {
	case "GET":
		related, _ := db.Communities_Related(id)
		sponsors, _ := db.Community_Sponsors(id)
		info, _ := db.Community_Info(id)
		info.Sponsors = sponsors
		info.Communities = related

		isMember, _ := db.Community_Membership_IsMember(id, username.(string))

		data := map[string]interface{}{
			"IsMember":  isMember,
			"Community": info,
		}

		// result := db.PRContributionAreas_Select()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		jsonResp, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonResp)

	case "POST":
		db.Community_Onboarding_AddMember(id, username.(string))
	case "DELETE":
		db.Community_Onboarding_RemoveMember(id, username.(string))
	}

}
