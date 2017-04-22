package models

import (
	"errors"

	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

const (
	// LogIndexPrefix ElasticSearch index name prefix
	LogIndexPrefix = "logs"
)

type Project struct {
	BaseModel
	Name    string `json:"name"`
	UUID    string `json:"uuid"`
	OwnerID uint   `json:"owner_id"`
}

var (
	ErrorProjectNotFound = errors.New("Project not found")
)

func CreateProject(name string, u *User) (*Project, *gorm.DB) {
	p := &Project{
		Name:    name,
		UUID:    uuid.NewV4().String(),
		OwnerID: u.ID,
	}
	res := db.Create(p)
	return p, res
}

// IndexName build ElasticSearch index name for project logs.
func (p *Project) IndexName() string {
	return ProjectIndexName(p.UUID)
}

func SaveProject(p *Project) (*Project, *gorm.DB) {
	res := db.Model(p).Save(p)
	return p, res
}

func FindProjectByUUID(uuid string) (*Project, *gorm.DB) {
	var p Project
	res := db.Model(&p).Where("uuid = ?", uuid).First(&p)
	return &p, res
}

func FindProjectById(id int, u *User) (*Project, *gorm.DB) {
	var p Project
	res := db.Model(&p).Where("id = ? AND owner_id = ?", id, u.ID).First(&p)
	return &p, res
}

func FindAllProjectsByUser(u *User) ([]Project, *gorm.DB) {
	var p []Project
	res := db.Model(&Project{}).Where("owner_id = ?", u.ID).Find(&p)
	return p, res
}

func DeleteProject(p *Project) *gorm.DB {
	return db.Delete(p)
}

// ProjectIndexName build ElasticSearch index name for project logs.
func ProjectIndexName(projectUUID string) string {
	return fmt.Sprintf("%s-%s", LogIndexPrefix, projectUUID)
}
