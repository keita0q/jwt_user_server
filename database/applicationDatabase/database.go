package applicationDatabase

import (
	"github.com/keita0q/user_server/model"
	"github.com/keita0q/user_server/database"
)

type Database interface {
	LoadChild(aChildID string) (*model.Child, *database.NotFoundError)
	LoadParent(aParentID string) (*model.Parent, *database.NotFoundError)
	LoadProject(aProjectID string) (*model.Project, *database.NotFoundError)
	SaveChild(*model.Child) error
	SaveParent(*model.Parent) error
	SaveProject(*model.Project) error
}

