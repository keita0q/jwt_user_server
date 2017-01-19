package sql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/keita0q/user_server/model"
	"github.com/keita0q/user_server/database"
)

type SqlDatabase struct {
	*gorm.DB
}

type Config struct {
	User      string
	ProjectID string
	DBName    string
}

func New() (*SqlDatabase, error) {
	tDB, tError := gorm.Open("mysql", "root:hohoho@/test")
	//tDB, tError := sql.Open("mysql", "user@cloudsql(project-id:instance-name)/dbname")
	//tDB, tError := sql.Open("mysql", "root@/test")
	if tError != nil {
		return nil, tError
	}
	return &SqlDatabase{tDB}, nil
}

func (aDB *SqlDatabase)LoadChild(aChildID string) (*model.Child, *database.NotFoundError) {
	//tRows, tError := aDB.db.Query("SELECT * FROM children")
	tChild := &model.Child{}
	aDB.Find(tChild, aChildID)
	return tChild, nil
}

func (aDB *SqlDatabase)LoadParent(aParentID string) (*model.Parent, *database.NotFoundError) {
	tParent := &model.Parent{}
	aDB.Find(tParent, aParentID)
	return tParent, nil
}

func (aDB *SqlDatabase)LoadProject(aProjectID string) (*model.Project, *database.NotFoundError) {
	tProject := &model.Project{}
	aDB.Find(tProject, aProjectID)
	return tProject, nil
}

func (aDB *SqlDatabase)SaveChild(aChild *model.Child) error {
	aDB.NewRecord(*aChild)
	aDB.Create(aChild)
	aDB.Save(aChild)
	return nil
}

func (aDB *SqlDatabase)SaveParent(aParent *model.Parent) error {
	aDB.NewRecord(*aParent)
	aDB.Create(aParent)
	aDB.Save(aParent)
	return nil
}

func (aDB *SqlDatabase)SaveProject(aProject *model.Project) error {
	aDB.NewRecord(*aProject)
	aDB.Create(aProject)
	aDB.Save(aProject)
	return nil
}
