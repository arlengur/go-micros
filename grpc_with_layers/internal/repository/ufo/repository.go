package ufo

import (
	"sync"

	def "github.com/arlengur/go-micros/layers/internal/repository"
	repoModel "github.com/arlengur/go-micros/layers/internal/repository/model"
)

var _ def.UFORepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Sighting
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]repoModel.Sighting),
	}
}
