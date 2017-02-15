package mail

type Mail interface {
	Send(aRcpt string, aMessage string) error
	GetDial() string
}
