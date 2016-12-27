package util

import (
	"net/http"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	from    = mail.NewEmail("Heroku Alert", os.Getenv("SENDGRID_USERNAME"))
	subject = "gops saw an error"
	to      = mail.NewEmail("ke", os.Getenv("SENDGRID_TARGET_ADDRESS"))
	content = mail.NewContent("text/plain", "")
)

func SendAlert() {
	m := mail.NewV3MailInit(from, subject, to, content)
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = http.MethodPost
	request.Body = mail.GetRequestBody(m)
	sendgrid.API(request)
}
