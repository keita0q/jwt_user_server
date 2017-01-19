package demo

import (
	"github.com/keita0q/user_server/database"
	"github.com/keita0q/user_server/model"
)

type DemoDB struct {
	children []model.Child
	parents  []model.Parent
	projects []model.Project
}

func New() *DemoDB {
	return &DemoDB{
		children: []model.Child{},
		parents: []model.Parent{},
		projects: []model.Project{},
	}
}

func (aDB *DemoDB)LoadChild(aChildID string) (*model.Child, *database.NotFoundError) {
	for _, tChild := range aDB.children {
		if tChild.ID == aChildID {
			return &tChild, nil
		}
	}
	return nil, &database.NotFoundError{Message:"notfound : " + aChildID}
}

func (aDB *DemoDB)LoadParent(aParentID string) (*model.Parent, *database.NotFoundError) {
	for _, tParent := range aDB.parents {
		if tParent.ID == aParentID {
			return &tParent, nil
		}
	}
	return nil, &database.NotFoundError{Message:"notfound : " + aParentID}
}

func (aDB *DemoDB)LoadProject(aProjectID string) (*model.Project, *database.NotFoundError) {
	for _, tProject := range aDB.projects {
		if tProject.ID == aProjectID {
			return &tProject, nil
		}
	}
	return nil, &database.NotFoundError{Message:"notfound : " + aProjectID}
}

func (aDB *DemoDB)SaveChild(aChild *model.Child) error {
	aDB.children = append(aDB.children, *aChild)
	return nil
}

func (aDB *DemoDB)SaveParent(aParent *model.Parent) error {
	aDB.parents = append(aDB.parents, *aParent)
	return nil
}

func (aDB *DemoDB)SaveProject(aProject *model.Project) error {
	aDB.projects = append(aDB.projects, *aProject)
	return nil
}
