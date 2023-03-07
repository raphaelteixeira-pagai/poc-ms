package uid

import "github.com/google/uuid"

type UIDGenerator interface {
	New() string
}

type uid struct{}

func (u uid) New() string {
	return uuid.New().String()
}
