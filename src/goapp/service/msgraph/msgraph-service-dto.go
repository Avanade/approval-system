package msgraph

type User struct {
	Name  string `json:"displayName"`
	Email string `json:"mail"`
}
