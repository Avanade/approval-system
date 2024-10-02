package msgraph

type MsGraphService interface {
	SearchUsers(search string) ([]User, error)
}
