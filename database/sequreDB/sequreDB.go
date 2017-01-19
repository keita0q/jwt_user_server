package sequreDB

import "github.com/keita0q/user_server/model"

type SequreDB interface {
	SaveUser(aUser *model.User) error
	LoadUser(aID string) (*model.User, error)
	Exist(aID string, aPassword string) bool
}
