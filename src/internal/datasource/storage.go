package datasource

import "sync"

type Storage struct {
	m sync.Map
}

func NewStorage() *Storage {
	return &Storage{}
}
