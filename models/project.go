package models

import (
	"errors"

	"github.com/satori/go.uuid"
)

type Project struct {
	ID      string
	Name    string
	UUID    string
	OwnerID int64
	Owner   *User
}

var (
	ErrorProjectNotFound = errors.New("Project not found")
)

func CreateProject(name string, u *User) error {
	p := &Project{
		Name:    name,
		UUID:    uuid.NewV4().String(),
		OwnerID: u.Id,
	}
	return db.Insert(p)
}

func FindProjectByUUID(uuid string) (*Project, error) {
	var p Project
	err := db.Model(&p).Where("uuid = ?", uuid).Select(&p)
	if err != nil {
		return nil, ErrorProjectNotFound
	}
	return &p, nil
}
