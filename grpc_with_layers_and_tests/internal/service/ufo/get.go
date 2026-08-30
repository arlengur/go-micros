package ufo

import (
	"context"

	"github.com/arlengur/go-micros/grpc_with_layers_and_tests/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	sighting, err := s.ufoRepository.Get(ctx, uuid)
	if err != nil {
		return model.Sighting{}, err
	}

	return sighting, nil
}
