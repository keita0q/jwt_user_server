package demoDB

import (
	"github.com/keita0q/user_server/model"
	"github.com/keita0q/user_server/database"
)

type DemoDB struct {
	users []model.User
}

func New() *DemoDB {
	return &DemoDB{
		users: []model.User{},
	}
}

func (aDB *DemoDB)SaveUser(aUser *model.User) error {
	aDB.users = append(aDB.users, *aUser)
	return nil
}

func (aDB *DemoDB)LoadUser(aID string) (*model.User, error) {
	for _, tUser := range aDB.users {
		if tUser.ID == aID {
			return &tUser, nil
		}
	}
	return nil, &database.NotFoundError{Message:"notfound : " + aID}
}

func (aDB *DemoDB)Exist(aID string, aPassword string) bool {
	for _, tUser := range aDB.users {
		if tUser.ID == aID {
			return true
		}
	}
	return false
}