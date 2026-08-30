package ufo

import (
	"github.com/arlengur/go-micros/grpc_with_layers_and_tests/internal/repository"
	def "github.com/arlengur/go-micros/grpc_with_layers_and_tests/internal/service"
)

var _ def.UFOService = (*service)(nil)

type service struct {
	ufoRepository repository.UFORepository
}

func NewService(ufoRepository repository.UFORepository) *service {
	return &service{
		ufoRepository: ufoRepository,
	}
}
