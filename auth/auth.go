package auth

type Auth interface {
	CreateToken(aID string, aPassword string) (string, error)
	Authenticate(aToken string) (Claim, bool, error)
}

type Claim interface {
	GetUserID() string
}

type NotFoundError struct {
	Message string
}

func (tError *NotFoundError) Error() string {
	return tError.Message
}