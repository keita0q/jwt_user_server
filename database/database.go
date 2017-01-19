package database

type NotFoundError struct {
	Message string
}

func (tError *NotFoundError) Error() string {
	return tError.Message
}