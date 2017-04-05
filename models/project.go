package models

import (
	"errors"

	"github.com/satori/go.uuid"
)

type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UUID      string `json:"uuid"`
	OwnerID   int64  `json:"owner_id"`
	CreatedAt string `json:"created_at"`
}

var (
	ErrorProjectNotFound = errors.New("Project not found")
)

func CreateProject(name string, u *User) (*Project, error) {
	p := &Project{
		Name:    name,
		UUID:    uuid.NewV4().String(),
		OwnerID: u.Id,
	}
	err := db.Insert(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func FindProjectByUUID(uuid string) (*Project, error) {
	var p Project
	err := db.Model(&p).Where("uuid = ?", uuid).Select()
	if err != nil {
		return nil, ErrorProjectNotFound
	}
	return &p, nil
}

func FindAllProjectsByUser(u *User) ([]Project, error) {
	var p []Project
	err := db.Model(&p).Where("owner_id = ?", u.Id).Select(&p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
