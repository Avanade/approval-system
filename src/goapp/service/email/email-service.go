package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"main/config"
	"main/model"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	graphusers "github.com/microsoftgraph/msgraph-sdk-go/users"
)

type sdkEmailService struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	UserId       string
	IsEnabled    bool
}

func NewSdkEmailService(config config.ConfigManager) EmailService {
	return &sdkEmailService{
		TenantID:     config.GetEmailTenantID(),
		ClientID:     config.GetEmailClientID(),
		ClientSecret: config.GetEmailClientSecret(),
		UserId:       config.GetEmailUserID(),
		IsEnabled:    config.GetIsEmailEnabled(),
	}
}

func (s *sdkEmailService) SendApprovalRequestEmail(req *model.ItemInsertRequest, appModule *model.ApplicationModule, id string) error {
	data := model.ApprovalRequestEmailData{
		Subject:     req.Subject,
		Body:        template.HTML(req.Body),
		ApproveText: appModule.ApproveText,
		RejectText:  appModule.RejectText,
		ApproveUrl:  fmt.Sprintf("%s/response/%s/%s/%s/1", os.Getenv("HOME_URL"), req.ApplicationId, req.ApplicationModuleId, id),
		RejectUrl:   fmt.Sprintf("%s/response/%s/%s/%s/0", os.Getenv("HOME_URL"), req.ApplicationId, req.ApplicationModuleId, id),
	}

	t, err := template.ParseFiles("templates/email.html", "templates/buttons.html")
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	msg := buf.String()

	err = s.SendEmail(req.Emails, nil, req.Subject, msg, Html, true)
	if err != nil {
		return err
	}

	return nil
}

func (s *sdkEmailService) SendActivityEmail(req *model.ItemActivity, recipients []string, domain, action string) error {
	bodyTemplate := ` 
			<tr style="color: #5c5c5c;"  >
				<td class="center-table" align="center">
					<table style="width: 100%; max-width: 700px;" class="margin-auto">
						<tr>
							<td style="padding: 15px 0 10px 0;">
								<span><b>|User|</b></span> 
								<span>commented on your request:</span>
							</td>
						</tr>
						<tr>
							<td align="center" style="padding: 10px;" bgcolor="#F8F8F8">
								<p>|Comment|</p>
							</td>
						</tr>
						<tr>
							<td style="padding: 15px 0;">
								<a href="|Domain|/|Action|/|AppId|/|AppModuleId|/|ItemId|/1?view=activities">
									View Conversation
								</a>
							</td>
						</tr>
					</table>
				</td>
            </tr>
            `

	replacer := strings.NewReplacer(
		"|User|", req.CreatedBy,
		"|Comment|", req.Content,
		"|AppId|", req.AppId,
		"|AppModuleId|", req.AppModuleId,
		"|ItemId|", req.ItemId,
		"|Domain|", domain,
		"|Action|", action,
	)

	htmlBody := s.buildHtmlBody(bodyTemplate, replacer)

	err := s.SendEmail(recipients, nil, "New Comment Notification", htmlBody, Html, false)
	if err != nil {
		return err
	}

	return nil
}

func (s *sdkEmailService) SendIPDRResponseEmail(data *model.IPDRequest, item *model.Item, domain string) error {
	action := "rejected"
	if data.IsApproved {
		action = "approved"
	}

	bodyTemplate := `
		<tr style="color: #5c5c5c;"  >
			<td class="center-table" align="center">
				<table style="width: 100%; max-width: 700px;" class="margin-auto">
					<tr>
						<td style="padding: 15px 0 0px 0;">
							<span>Hi,</span> 
						</td>
					</tr>
					<tr>
						<td style="padding: 10px 0 10px 0;">
							<span>The following request for Intellectual Property disclosure has been |Action|.</span>
						</td>
					</tr>
					<tr>
						<td style="padding: 10px 0 10px 0;">
							<span>Approver Remarks: |Remarks|.</span>
						</td>
					</tr>
					<tr>
						<td class="center-table">
							<table style="width: 100%; max-width: 700px; margin: auto">
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Requestor
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Requestor|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Requestor Email
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Email|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Involvement
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Involvement|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Intellectual Property Title
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|IPTitle|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Intellectual Property Type
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|IPType|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Intellectual Property Description
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|IPDescription|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Reason
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Reason|
									</td>
								</tr>
							</table>
						</td>
					</tr>
				</table>
			</td>
		</tr>
		`

	replacer := strings.NewReplacer(
		"|Action|", action,
		"|Remarks|", data.ApproverRemarks,
		"|Requestor|", data.RequestorName,
		"|Email|", data.RequestorEmail,
		"|Involvement|", strings.Join(data.Involvement, ", "),
		"|IPTitle|", data.IPTitle,
		"|IPType|", data.IPType,
		"|IPDescription|", data.IPDescription,
		"|Reason|", data.Reason,
	)
	htmlBody := s.buildHtmlBody(bodyTemplate, replacer)

	err := s.SendEmail([]string{data.RequestorEmail}, nil, "IP Disclosure Request", htmlBody, Html, false)
	if err != nil {
		return err
	}

	return nil
}

func (s *sdkEmailService) SendLegalConsultationRequestEmail(req *model.ConsultLegalRequest, user *model.AzureUser, domain string, recipients []string) error {
	bodyTempate := ` 
			<tr style="color: #5c5c5c;"  >
				<td class="center-table" align="center">
					<table style="width: 100%; max-width: 700px;" class="margin-auto">
						<tr>
							<td style="padding: 15px 0 0px 0;">
								<span>Hi,</span> 
							</td>
						</tr>
						<tr>
							<td style="padding: 10px 0 10px 0;">
								<span><a href="mailto:|UserEmail|">|User|</a></span> 
								<span>is requesting for legal review and your input in an IP disclosure request.</span>
							</td>
						</tr>
						<tr>
							<td style="padding: 15px 0;">
								<a href="|Domain|/review/|AppId|/|AppModuleId|/|ItemId|/1">
									View Request
								</a>
							</td>
						</tr>
					</table>
				</td>
            </tr>
            `

	replacer := strings.NewReplacer(
		"|User|", user.Name,
		"|UserEmail|", user.Email,
		"|AppId|", req.ApplicationId,
		"|AppModuleId|", req.ApplicationModuleId,
		"|ItemId|", req.ItemId,
		"|Domain|", domain,
	)

	htmlBody := s.buildHtmlBody(bodyTempate, replacer)

	err := s.SendEmail(recipients, nil, "IP Disclosure Request", htmlBody, Html, false)
	if err != nil {
		return err
	}

	return nil
}

func (s *sdkEmailService) SendEmail(to, cc []string, subject, content string, contentType ContentType, isSaveToSentItem bool) error {
	requestBody := graphusers.NewItemSendMailPostRequestBody()
	message := graphmodels.NewMessage()
	message.SetSubject(&subject)
	body := graphmodels.NewItemBody()
	bodyType := graphmodels.TEXT_BODYTYPE
	if contentType == Html {
		bodyType = graphmodels.HTML_BODYTYPE
	}
	body.SetContentType(&bodyType)
	body.SetContent(&content)
	message.SetBody(body)

	var toRecipients []graphmodels.Recipientable
	for _, v := range to {
		recipient := graphmodels.NewRecipient()
		emailAddress := graphmodels.NewEmailAddress()
		address := v
		emailAddress.SetAddress(&address)
		recipient.SetEmailAddress(emailAddress)
		toRecipients = append(toRecipients, recipient)
	}
	message.SetToRecipients(toRecipients)

	var ccRecipients []graphmodels.Recipientable
	for _, v := range cc {
		recipient := graphmodels.NewRecipient()
		emailAddress := graphmodels.NewEmailAddress()
		address := v
		emailAddress.SetAddress(&address)
		recipient.SetEmailAddress(emailAddress)
		ccRecipients = append(ccRecipients, recipient)
	}
	message.SetCcRecipients(ccRecipients)

	requestBody.SetMessage(message)
	requestBody.SetSaveToSentItems(&isSaveToSentItem)

	graphServiceClient, err := s.connectGraphServiceClient()
	if err != nil {
		return err
	}

	return graphServiceClient.Users().ByUserId(s.UserId).SendMail().Post(context.Background(), requestBody, nil)
}

func (s *sdkEmailService) connectGraphServiceClient() (*msgraphsdk.GraphServiceClient, error) {
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

func (s *sdkEmailService) headStyle() string {
	return `<head>
				<style>
					table, th, tr, td {
						border: 0;
						border-collapse: collapse;
						vertical-align: middle;
					}
					.thead {
						padding: 15px;
					}
					.center-table {
						text-align: -webkit-center;
					}
					.margin-auto {
						margin: auto;
					}
					.border-top {
						border-top: 1px rgb(204, 204, 204) solid;
						border-collapse: separate;
					}
				</style>
			</head>`
}

func (s *sdkEmailService) header() string {
	return `<tr>
                <th class="center-table">
                    <table style="width: 100%; max-width: 700px;" class="margin-auto">
                        <tr style="background-color: #ff5800">
                            <td class="thead" style="width: 95px">
                                <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGMAAAAdCAQAAAAUGhqvAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QA/4ePzL8AAAAHdElNRQfmCQcLLziqRWflAAAGWUlEQVRYw92YaXCV5RXHf29yL0kg7IRFEAqGDOAwrEKxBBzABfrFjowsUingCFIXFrF0mREHtVLb6ahdLINKW6WKChRsKVsbgqxJS1gCsUhlK8giS4Ds9/76ITc3NwH80C8Yz/3w3vd/zjvv8z/nPP/neV74WlhQHxACWpDCeSqu424YJjZzrrnu8k17ijd7QP8niXRfdZsTHeFS/+FtDZCGiLMtsI/VVfmjS23c4IiIfd3reIn9epjvhAbWWGLYxb5lijX3+JQ5tm1ANEQcYoF31GZfbGuuMxtQPcSQr/uaySZi+Ihb7GBrs2xlcDPpiIE9zbbLDdMq4kDzHGB9vI0bfdccd5vjQ4ZvHhEx5HtGXHA9GkmxazJT2E7BNf7zFHE/q3mY95jPBG5mgwWkkET4eq5Q7DqYO5hOJHHVFuBe7uI4X7CPfVxmNnkcvGk0vmSShkBowhNsYs813s78hHc4wzQ2cJIPuI+pPKNBLc2OjCCLJnzBTrZSSncGU8JGigOEzgyllI00ZyTdacw5drCNMiCZLvShOxnAcXIoJAKkks0AMqjgIJv4LwQIIfqSTUfKuC0hxZ25myzK2UkOV6vnxSS32uWaeRHyl66wuU1d4UKTxZHusGuCkuHT1thlXzJsthe84n2xlefHao7pzo9HFfuCqWInd1kZR485UcSenoghEbfYS8Q0f+hpa+0FEb9lvlGrjFriq7ZA7OdOZ9vXUY6yj80MYoMY7e5qAXaw+T4opvsXp9XUVsRRLnSqU33bCosdYRP/ri6yemOzSZ0v3uNzTnGay6z0osPFTE8Y9U8+4Xx3qfu9VWzvsz7ud33Wo+prJonTLbXUPzjFxzwYo9Herer7ftsZHrHCOTjIzZ5wnZvd4nYL3OjTfkPMcIPz4yv6eP/pw4Z90rdr9UoMYqRb+rE6T5yrbrel2N9znrFfQlQbt6uzYjTKHSVithcs9e6EuOoqb7OZ7fyX+hvTxJArYzQeMmphtfQ6T80PsYghFLCF9XxKEu0Yxv2M411aEWFJ9VZdWE45c7iLIjK5haPx3gszmFHcSoR2QDNgA59zO73YSjatWcVBIMwQRtKJCBlA87oqwhHOkkkLQLoymt6E6Qg0JZkssrjCckqBIH5uGEBACfcQJUonImSGOEAjNjCI0XzEEgop5B3uZTaDWEU3IhYTAaKsopCpPEIzOlTTEEI8yY9oSQlVNKGa8yd8zFiGsYuRwBrKCDOb+bSoE5VoESIxdBC/pT+VlJASQzJI5Qyn6j2RAQxkYPw+NcQH9GEZLzOCWWQzj0KK+RvjWEsVi7nMCc5TToimtCCdVnxIYfzx3swlnZ+zkip+yggAylnDWIazgQH8hxygD3NpzCL+TJSfMeyGcprEHPqTy8scZwwvAlCJhEmrF1sG5LEkLsCREDv4jOn8gDXksZDXeZQiZtKeSZykG73IpDWNKOckxxhJKQu4HE9nJu04xCucIMS5+EtyOURfJtOBNzgKdKcNn/AKpwgnRF1Lozk9gKV8BHRHIOAYF2jDYArqrBtFgKzkbM1IQpTyKxbzHd7nc+byIr9gIxN4hiMQFFGU0MQT6cn36xT4ImVkcCfraUTjOHqMTcxgMuWsIQJcpIK23MkmUhKi6ltAOZeAb7KeS7QkAJL5N7k8wDxKySNEy1jsBo4zkIUs8Txp3MKBarWZ5l4nmS5mWWB5/V1tTHR3Oy4RF1u5Wi12n/u9qj4f05kxlqj5to3p01/VS/GoBfWUqqMH1QfFx7xqlYfd7Sl1vxliP/PUSk972sqYUiX5qOfUix7zrMVOCAUov6eKWUzmDF24zB56k86VOrlqx3Os48N68/M8j5NPNhlUsY1DrI3hO/g1nVjPWQDOMZPvMZQ2VLGVQ6wDLrOCprHKlrCaPI4Cb3GeB+hKmE/ZTA5XgN2MYwJDaUuUPRxmLRDlTQ4znttJpZi9FAaxlgnoxnDa8hmbac0brOd5KuKbjka8RA8m1/ZibUUISCMNKaUMgxos5g7qRkUpoyzW9bXdmvgfwqSTRAWlVMXlPiCFJonvACGZdEKUU0L0OjNNHOZe5xgSappujwMb0PEpTmSMe51lqpjsWAur9zwNzGJEdrvMp/ydB5xRXZmvrt3gs6AQ0JtJZHKa5eTWPYl89exLRicEJBMl+tX/CPo/520riCgLgNcAAAAldEVYdGRhdGU6Y3JlYXRlADIwMjItMDktMDdUMTE6NDc6NTYrMDA6MDA42qGMAAAAJXRFWHRkYXRlOm1vZGlmeQAyMDIyLTA5LTA3VDExOjQ3OjU2KzAwOjAwSYcZMAAAAABJRU5ErkJggg==" />
                            </td>
                            <td class="thead" style="font-family: SegoeUI, sans-serif;font-size: 14px;color: white;padding-bottom: 10px;">
                                Community
                            </td>
                        </tr>
                    </table>
                </th>
            </tr>`
}

func (s *sdkEmailService) buildHtmlBody(bodyTemplate string, replacer *strings.Replacer) string {
	replacedBodyTemplate := replacer.Replace(bodyTemplate)
	htmlBody := `<html>`
	htmlBody += s.headStyle()
	htmlBody += `<body>`
	htmlBody += `<table style="width: 100%">`
	htmlBody += s.header()
	htmlBody += replacedBodyTemplate
	htmlBody += `</table>`
	htmlBody += `</body>`
	htmlBody += `</html>`
	return htmlBody
}
