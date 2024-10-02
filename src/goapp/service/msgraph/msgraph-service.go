package msgraph

import (
	"context"
	"fmt"
	"main/config"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

type msGraphService struct {
	TenantID     string
	ClientID     string
	ClientSecret string
}

func NewMsGraphService(config config.ConfigManager) MsGraphService {
	return &msGraphService{
		TenantID:     config.GetTenantID(),
		ClientID:     config.GetClientID(),
		ClientSecret: config.GetClientSecret(),
	}
}

func (service *msGraphService) SearchUsers(search string) ([]User, error) {
	graphServiceClient, err := service.connectGraphServiceClient()
	if err != nil {
		return nil, err
	}

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	filter := fmt.Sprintf(`"displayName:%s" OR "mail:%s" OR "otherMails:%s"`, search, search, search)

	fields := users.UsersRequestBuilderGetQueryParameters{
		Select: []string{"displayName", "otherMails", "mail"},
		Search: &filter,
	}

	options := users.UsersRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: &fields,
	}

	users, err := graphServiceClient.Users().Get(context.Background(), &options)
	if err != nil {
		return nil, err
	}

	var result []User
	for _, user := range users.GetValue() {
		name := user.GetDisplayName()
		email := user.GetMail()
		otherEmail := user.GetOtherMails()

		u := User{
			Name: *name,
		}

		if *email != "" {
			u.Email = *email
		} else if len(otherEmail) > 0 {
			u.Email = otherEmail[0]
		}

		result = append(result, u)
	}

	return result, nil
}

func (s *msGraphService) connectGraphServiceClient() (*msgraphsdk.GraphServiceClient, error) {
	cred, _ := azidentity.NewClientSecretCredential(
		s.TenantID,
		s.ClientID,
		s.ClientSecret,
		nil,
	)

	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(
		cred, []string{"https://graph.microsoft.com/.default"},
	)

	if err != nil {
		return nil, err
	}

	return graphClient, nil
}
