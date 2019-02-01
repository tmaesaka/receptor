package store

import (
	"errors"
)

// Store is an interface that Storage Engines must implement.
type Store interface {
	Name() string
}

// GetStore returns a Storage Engine based on the given name.
func GetStore(name string) (Store, error) {
	switch name {
	case "memory":
		return MemoryStore{}, nil
	case "mysql", "mariadb":
		return nil, errors.New("work in progress")
	default:
		return nil, errors.New("unknown store")
	}
}
