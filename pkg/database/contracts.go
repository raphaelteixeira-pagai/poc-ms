package database

import "github.com/gocraft/dbr/v2"

type DBPool interface {
	Acquire() (*dbr.Session, error)
	Release(*dbr.Session)
	Close()
}
