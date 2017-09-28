package smtp

import (
	"net/smtp"
)

type SmtpMail struct {
	address string
	auth    smtp.Auth
	dial    string
}

type Config struct {
	Address  string
	NeedAuth bool
	Identity string
	UserName string
	Password string
	Host     string
}

func New(aConfig *Config) *SmtpMail {
	if aConfig.NeedAuth {
		tAuth := smtp.PlainAuth(aConfig.Identity, aConfig.UserName, aConfig.Password, aConfig.Host)
		return &SmtpMail{address: aConfig.Address, auth: tAuth, dial: aConfig.UserName}
	}
	return &SmtpMail{address: aConfig.Address, auth: nil, dial: aConfig.UserName}
}

func (aMail *SmtpMail) Send(aRcpt string, aMessage string) error {
	return smtp.SendMail(
		aMail.address,
		aMail.auth,
		aMail.dial,
		[]string{aRcpt},
		[]byte(aMessage),
	)
}

func (aMail *SmtpMail) GetDial() string {
	return aMail.dial
}
