package models

import (
	"errors"

	"github.com/satori/go.uuid"
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

func CreateProject(name string, u *User) (*Project, error) {
	p := &Project{
		Name:    name,
		UUID:    uuid.NewV4().String(),
		OwnerID: u.ID,
	}
	res := db.Create(p)
	if res.Error != nil {
		return nil, res.Error
	}
	return p, nil
}

func FindProjectByUUID(uuid string) (*Project, error) {
	var p Project
	res := db.Model(&p).Where("uuid = ?", uuid).First(&p)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RecordNotFound() {
		return nil, ErrorProjectNotFound
	}
	return &p, nil
}

func FindProjectById(id int, u *User) (*Project, error) {
	var p Project
	res := db.Model(&p).Where("id = ? AND owner_id = ?", id, u.ID).First(&p)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RecordNotFound() {
		return nil, ErrorProjectNotFound
	}
	return &p, nil
}

func FindAllProjectsByUser(u *User) ([]Project, error) {
	var p []Project
	res := db.Model(&Project{}).Where("owner_id = ?", u.ID).Find(&p)
	if res.Error != nil {
		return nil, res.Error
	}
	return p, nil
}
