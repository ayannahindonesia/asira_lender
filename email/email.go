package email

import (
	"asira_lender/asira"
	"fmt"

	"gopkg.in/gomail.v2"
)

func SendMail(to string, subject, message string) error {
	Config := asira.App.Config.GetStringMap(fmt.Sprintf("%s.mailer", asira.App.ENV))
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", Config["email"].(string))
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", message)

	dialer := gomail.NewPlainDialer(Config["host"].(string),
		Config["port"].(int),
		Config["email"].(string),
		Config["password"].(string))

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}
	//

	return nil
}
